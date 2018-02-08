package bhattacharya

import (
	"encoding/csv"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync/atomic"

	"go.uber.org/zap"

	"core/utils"
)

const defaultProb = 0.00000000001

var ErrUnderflow = errors.New("possible underflow detected")

type NBClass string

// NBClassifier implements the Naive Bayesian Classifier.
type NBClassifier struct {
	Classes         []NBClass
	learned         int
	seen            int32
	datas           map[NBClass]*classData
	rawDatas        map[NBClass]*classData
	tfIdf           bool
	DidConvertTfIdf bool
}

// DOC: serializableClassifier represents a container for Classifier objects
//      whose fields are modifiable by reflection and are writeable by GOB.
type serializableClassifier struct {
	Classes         []NBClass
	Learned         int
	Seen            int
	Datas           map[NBClass]*classData
	RawDatas        map[NBClass]*classData
	TfIdf           bool
	DidConvertTfIdf bool
}

// DOC: classData holds the frequency data for words in a particular class.
type classData struct {
	Freqs   map[string]float64
	FreqTfs map[string][]float64
	Total   int
}

func newClassData() *classData {
	return &classData{
		Freqs:   make(map[string]float64),
		FreqTfs: make(map[string][]float64),
	}
}

// DOC: The probability of seeing a word in a document of this class.
func (d *classData) getWordProb(word string) float64 {
	value, ok := d.Freqs[word]
	if !ok {
		return defaultProb
	}
	return float64(value) / float64(d.Total)
}

// DOC: The probability of seeing a set of words in a document of this class.
func (d *classData) getWordsProb(words []string) (prob float64) {
	prob = 1
	for _, word := range words {
		prob *= d.getWordProb(word)
	}
	return
}

func classifierChecks(classes ...NBClass) {
	n := len(classes)
	if n < 2 {
		utils.ModelLog.Panic("Provide at least two classes.")
	}

	check := make(map[NBClass]bool, n)
	for _, class := range classes {
		check[class] = true
	}
	if len(check) != n {
		utils.ModelLog.Panic("Model classes must be unique.")
	}
}

func NewNBClassifierTfIdf(classes ...NBClass) (c *NBClassifier) {
	classifierChecks(classes...)
	c = &NBClassifier{
		Classes:  classes,
		datas:    make(map[NBClass]*classData, len(classes)),
		rawDatas: make(map[NBClass]*classData, len(classes)),
		tfIdf:    true,
	}
	for _, class := range classes {
		c.datas[class] = newClassData()
		c.rawDatas[class] = newClassData()
	}
	return
}

func NewNBClassifier(classes ...NBClass) (c *NBClassifier) {
	classifierChecks(classes...)
	c = &NBClassifier{
		Classes:         classes,
		datas:           make(map[NBClass]*classData, len(classes)),
		rawDatas:        make(map[NBClass]*classData, len(classes)),
		tfIdf:           false,
		DidConvertTfIdf: false,
	}
	for _, class := range classes {
		c.datas[class] = newClassData()
		c.rawDatas[class] = newClassData()
	}
	return
}

// DOC: NewClassifierFromFile loads an existing classifier from file.
func NewNBClassifierFromFile(name string) (c *NBClassifier, err error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return NewNBClassifierFromReader(file)
}

// DOC: This actually does the deserializing of a GOB encoded classifier
func NewNBClassifierFromReader(r io.Reader) (c *NBClassifier, err error) {
	dec := gob.NewDecoder(r)
	w := new(serializableClassifier)
	err = dec.Decode(w)

	return &NBClassifier{w.Classes, w.Learned, int32(w.Seen), w.Datas, w.RawDatas, w.TfIdf, w.DidConvertTfIdf}, err
}

func (c *NBClassifier) getPriors() (priors []float64) {
	n := len(c.Classes)
	priors = make([]float64, n, n)
	sum := 0
	smoother := 1.0
	for index, class := range c.Classes {
		total := c.datas[class].Total
		priors[index] = float64(total)
		sum += total
	}
	if sum != 0 {
		for i := 0; i < n; i++ {
			// DOC: This is Laplace smoothing implemented below.
			priors[i] += smoother
			priors[i] /= (float64(sum) + (smoother * float64(n)))
		}
	}
	return
}

func (c *NBClassifier) Learned() int {
	return c.learned
}

func (c *NBClassifier) Seen() int {
	return int(atomic.LoadInt32(&c.seen))
}

func (c *NBClassifier) IsTfIdf() bool {
	return c.tfIdf
}

// DOC: WordCount finds the number of words for each class in the classifier.
func (c *NBClassifier) WordCount() (result []int) {
	result = make([]int, len(c.Classes))
	for inx, class := range c.Classes {
		data := c.datas[class]
		result[inx] = data.Total
	}
	return
}

func (c *NBClassifier) Observe(word string, count int, which NBClass) {
	data := c.datas[which]
	data.Freqs[word] += float64(count)
	data.Total += count
}

// DOC: Learn will accept new training documents for supervised learning.
func (c *NBClassifier) Learn(document []string, which NBClass) {
	if c.tfIdf {
		if c.DidConvertTfIdf {
			//utils.ModelSummary.Panic("Cannot call ConvertTermsFreqToTfIdf more than once. Reset and relearn to reconvert.")
		}

		// DOC: Term frequency: word count in document / document length.
		docTf := make(map[string]float64)
		for _, word := range document {
			docTf[word]++
		}

		docLen := float64(len(document))
		for wIndex, wCount := range docTf {
			docTf[wIndex] = wCount / docLen
			c.datas[which].FreqTfs[wIndex] = append(c.datas[which].FreqTfs[wIndex], docTf[wIndex])
			c.rawDatas[which].FreqTfs[wIndex] = append(c.rawDatas[which].FreqTfs[wIndex], docTf[wIndex])
		}
	}

	data := c.datas[which]
	rawData := c.rawDatas[which]
	for _, word := range document {
		data.Freqs[word]++
		data.Total++
		rawData.Freqs[word]++
		rawData.Total++
	}
	c.learned++
}

func (c *NBClassifier) OnlineLearn(document []string, which NBClass) {
	if !c.tfIdf {
		c.Learn(document, which)
	}
	// DOC: Term frequency: word count in document / document length.
	docTf := make(map[string]float64)
	for _, word := range document {
		docTf[word]++
	}

	if c.rawDatas[which] == nil {
		c.rawDatas[which] = newClassData()
		c.datas[which] = newClassData()
	}

	docLen := float64(len(document))
	for wIndex, wCount := range docTf {
		docTf[wIndex] = wCount / docLen
		c.rawDatas[which].FreqTfs[wIndex] = append(c.rawDatas[which].FreqTfs[wIndex], docTf[wIndex])
		c.datas[which].FreqTfs[wIndex] = append(c.datas[which].FreqTfs[wIndex], docTf[wIndex])
	}

	data := c.datas[which]
	rawData := c.rawDatas[which]
	for _, word := range document {
		data.Freqs[word]++
		data.Total++
		rawData.Freqs[word]++
		rawData.Total++
	}
	c.learned++

	for className, _ := range c.rawDatas {
		sample := math.Log1p(float64(c.learned) / float64(c.rawDatas[className].Total))
		for wIndex, _ := range c.rawDatas[className].FreqTfs {
			tfIdfAdder := float64(0)
			freqTfs := c.datas[className].FreqTfs[wIndex]
			for tfSampleIndex, tf := range c.rawDatas[className].FreqTfs[wIndex] {
				result := math.Log1p(tf) * sample
				tfIdfAdder += result
				freqTfs[tfSampleIndex] = result
			}
			c.datas[className].Freqs[wIndex] = tfIdfAdder
		}
	}
}

func (c *NBClassifier) ConvertTermsFreqToTfIdf() {
	if c.DidConvertTfIdf {
		utils.ModelLog.Panic("Cannot call ConvertTermsFreqToTfIdf more than once. Reset and relearn to reconvert.")
	}
	for className, _ := range c.datas {
		for wIndex, _ := range c.datas[className].FreqTfs {
			tfIdfAdder := float64(0)
			for tfSampleIndex, _ := range c.datas[className].FreqTfs[wIndex] {
				tf := c.datas[className].FreqTfs[wIndex][tfSampleIndex]
				c.datas[className].FreqTfs[wIndex][tfSampleIndex] = math.Log1p(tf) * math.Log1p(float64(c.learned)/float64(c.datas[className].Total))
				tfIdfAdder += c.datas[className].FreqTfs[wIndex][tfSampleIndex]
			}
			c.datas[className].Freqs[wIndex] = tfIdfAdder
		}
	}
	c.DidConvertTfIdf = true
}

// DOC: LogScores produces "log-likelihood"-like scores that can be used to
//      classify documents into classes.
func (c *NBClassifier) LogScores(document []string) (scores []float64, inx int, strict bool) {
	if c.tfIdf && !c.DidConvertTfIdf {
		utils.ModelLog.Panic("Using a TF-IDF classifier. Please call ConvertTermsFreqToTfIdf before calling LogScores.")
	}

	n := len(c.Classes)
	scores = make([]float64, n, n)
	priors := c.getPriors()
	for index, class := range c.Classes {
		data := c.datas[class]
		score := math.Log(priors[index])
		for _, word := range document {
			score += math.Log(data.getWordProb(word))
		}
		scores[index] = score
	}
	inx, strict = findMax(scores)
	atomic.AddInt32(&c.seen, 1)
	return scores, inx, strict
}

// DOC: Delivers actual probabilities.
func (c *NBClassifier) ProbScores(doc []string) (scores []float64, inx int, strict bool) {
	if c.tfIdf && !c.DidConvertTfIdf {
		utils.ModelLog.Panic("Using a TF-IDF classifier. Please call ConvertTermsFreqToTfIdf before calling ProbScores.")
	}
	n := len(c.Classes)
	scores = make([]float64, n, n)
	priors := c.getPriors()
	sum := float64(0)
	for index, class := range c.Classes {
		data := c.datas[class]
		score := priors[index]
		for _, word := range doc {
			score *= data.getWordProb(word)
		}
		scores[index] = score
		sum += score
	}
	for i := 0; i < n; i++ {
		scores[i] /= sum
	}
	inx, strict = findMax(scores)
	atomic.AddInt32(&c.seen, 1)
	return scores, inx, strict
}

func (c *NBClassifier) SafeProbScores(doc []string) (scores []float64, inx int, strict bool, err error) {
	if c.tfIdf && !c.DidConvertTfIdf {
		utils.ModelLog.Panic("Using a TF-IDF classifier. Please call ConvertTermsFreqToTfIdf before calling SafeProbScores.")
	}

	n := len(c.Classes)
	scores = make([]float64, n, n)
	logScores := make([]float64, n, n)
	priors := c.getPriors()
	sum := float64(0)
	for index, class := range c.Classes {
		data := c.datas[class]
		score := priors[index]
		logScore := math.Log(priors[index])
		for _, word := range doc {
			p := data.getWordProb(word)
			score *= p
			logScore += math.Log(p)
		}
		scores[index] = score
		logScores[index] = logScore
		sum += score
	}
	for i := 0; i < n; i++ {
		scores[i] /= sum
	}
	inx, strict = findMax(scores)
	logInx, logStrict := findMax(logScores)

	// DOC: This detects underflow errors.
	if inx != logInx || strict != logStrict {
		err = ErrUnderflow
	}
	atomic.AddInt32(&c.seen, 1)
	return scores, inx, strict, err
}

func (c *NBClassifier) WordFrequencies(words []string) (freqMatrix [][]float64) {
	n, l := len(c.Classes), len(words)
	freqMatrix = make([][]float64, n)
	for i, _ := range freqMatrix {
		arr := make([]float64, l)
		data := c.datas[c.Classes[i]]
		for j, _ := range arr {
			arr[j] = data.getWordProb(words[j])
		}
		freqMatrix[i] = arr
	}
	return
}

// DOC: Returns a map of words and their probability for a given class.
func (c *NBClassifier) WordsByClass(class NBClass) (freqMap map[string]float64) {
	freqMap = make(map[string]float64)
	for word, cnt := range c.datas[class].Freqs {
		freqMap[word] = float64(cnt) / float64(c.datas[class].Total)
	}
	return freqMap
}

func (c *NBClassifier) WriteToFile(name string) (err error) {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return c.WriteTo(file)
}

func (c *NBClassifier) WriteClassesToFile(rootPath string) (err error) {
	for name, _ := range c.datas {
		c.WriteClassToFile(name, rootPath)
	}
	return
}

func (c *NBClassifier) WriteClassToFile(name NBClass, rootPath string) (err error) {
	data := c.datas[name]
	fileName := filepath.Join(rootPath, string(name))
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(file)
	err = enc.Encode(data)
	return
}

func (c *NBClassifier) WriteTo(w io.Writer) (err error) {
	enc := gob.NewEncoder(w)
	err = enc.Encode(&serializableClassifier{c.Classes, c.learned, int(c.seen), c.datas, c.rawDatas, c.tfIdf, c.DidConvertTfIdf})

	return
}

func (c *NBClassifier) ReadClassFromFile(class NBClass, location string) (err error) {
	fileName := filepath.Join(location, string(class))
	file, err := os.Open(fileName)

	if err != nil {
		return err
	}

	dec := gob.NewDecoder(file)
	w := new(classData)
	err = dec.Decode(w)

	c.learned++
	c.datas[class] = w
	return
}

func findMax(scores []float64) (inx int, strict bool) {
	inx = 0
	strict = true
	for i := 1; i < len(scores); i++ {
		if scores[inx] < scores[i] {
			inx = i
			strict = true
		} else if scores[inx] == scores[i] {
			strict = false
		}
	}
	return
}

func (c *NBClassifier) LogClassWords() {
	for class, words := range c.datas {
		utils.ModelLog.Info("---------", zap.String("Class", string(class)))
		type wordCount struct {
			word  string
			count int
		}
		temp := []wordCount{}
		for k, v := range words.Freqs {
			temp = append(temp, wordCount{k, int(v)})
		}
		sort.Slice(temp, func(i, j int) bool {
			return temp[i].count > temp[j].count
		})
		limit := 9
		if len(temp) < limit {
			limit = len(temp)
		}
		for _, kv := range temp[:limit] {
			utils.ModelLog.Info("", zap.String("Word", string(kv.word)), zap.Int("Count", int(kv.count)))
		}
	}
}

func (c *NBClassifier) GenerateProbabilityTable(issueID int64, content string, assignees []string, status string) {
	// set up empty matrix
	words := strings.Split(content, " ")
	results := make([][]string, (len(words) + 1))
	for i := range results {
		results[i] = make([]string, (len(assignees) + 1))
	}

	// fill out column headers
	results[0][0] = "words"
	copy(results[0][1:], assignees)

	// populate matrix values
	for i := range words {
		results[i+1][0] = words[i]
		for j := range assignees {
			data := c.datas[NBClass(assignees[j])]
			results[i+1][j+1] = fmt.Sprintf("%f", math.Log(data.getWordProb(words[i])))
		}
	}

	file, _ := os.Create(fmt.Sprintf("../../../../data/backtests/%v-results-issue-%d.csv", status, issueID))
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	for _, value := range results {
		writer.Write(value)
	}
}
