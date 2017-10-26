go build -tags prod
cp ./cmd /home/heupr/golang/bin/core/pipeline/backend/cmd/cmd
cp ./startprod.sh /home/heupr/golang/bin/core/pipeline/backend/cmd/start.sh
chmod +x /home/heupr/golang/bin/core/pipeline/backend/cmd/start.sh
cp ./configprod.yaml /home/heupr/golang/bin/core/pipeline/backend/cmd/config.yaml
cp ./heupr.2017-10-04.private-key.pem
