Actionable Items - "Replicating jss12bhattacharya"
0) Enter below actionable items into issue tracker and add skeleton code
1) Benchmark Training Data
   -Mozilla May 1998 - Present
   -Eclispe Octobor 2001 - Present
   "Moreover, our data sets are much larger, covering the entire lifespan of both Mozilla (from May 1998 to March 2010) and Eclipse (from October 2001 to March 2010)."
2) Download Core-FX Github Training Data
3) Extract Fields from Bug Report Training Data
   -Dev ID
   -Product Name
   -Product Component
   -Time it took to fix
   -Relevant Words (see 3 tf-idf)
   1. Keywords: we collect keywords from the bug title, bug description
      and comments in the bug report.
   2. Bug source: we retrieve the product and component the bug has
      been filed under from the bug report.
   3. Temporal information: we collect information about when the
      bug has been reported and when it has been fixed.
   4. Developers assigned: we collect the list of developer IDs assigned
      to the bug from the activity page of the bug and the bug routing
      sequence. "
4) Implement Tf-idf algorithm for extracting relevant words
   ". For extracting relevant words
   from bug reports, we employ tf-idf, stemming, stop-word and
   non-alphabetic word removal(Manning et al., 2008). We use the
   Weka toolkit(Weka Toolkit, 2010) to remove stop words and form
   the word vectors for the dictionary (via the StringtoWordVector
   class with tf-idf enabled)."
5) Implement Classifiers
   -Naive Bayes
   -Bayesian Networks
   -SVM (Polynomial and RBF)
6) Implement 10-Fold Cross Validation      
7) Implement Incremental learning
8) Implement Ranking Function For Tossing Graph
     Rank(Tk) = Pr(Di -> Tk) + MatchedProduct(Tk) + MatchedComponent(Tk) + LastActivity(Tk)
(New Research Ideas)
9) Implement/Investigate the use of SyntaxNet (Original research opportunity)
10) Implement/Investigate the use of Labels (If #9 is a grand slam than #10 is not needed)
