go build -tags prod
cp ./cmd /home/heupr/golang/bin/core/pipeline/ingestor/cmd/cmd
cp ./startprod.sh /home/heupr/golang/bin/core/pipeline/ingestor/cmd/start.sh
chmod +x /home/heupr/golang/bin/core/pipeline/ingestor/cmd/start.sh
cp ./configprod.yaml /home/heupr/golang/bin/core/pipeline/ingestor/cmd/config.yaml
cp ./heupr.2017-10-04.private-key.pem /home/heupr/golang/bin/core/pipeline/ingestor/cmd/heupr.2017-10-04.private-key.pem
