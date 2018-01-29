go build -tags prod
cp ./cmd /home/heupr/golang/bin/core/pipeline/frontend/cmd/cmd
cp ./startprod.sh /home/heupr/golang/bin/core/pipeline/frontend/cmd/start.sh
chmod +x /home/heupr/golang/bin/core/pipeline/frontend/cmd/start.sh
cp ./prodprep.sh /home/heupr/golang/bin/core/pipeline/frontend/cmd/prep.sh
cp ./configprod.yaml /home/heupr/golang/bin/core/pipeline/frontend/cmd/config.yaml
cp -r ../static/ /home/heupr/golang/bin/core/pipeline/frontend
cp -r ../templates/ /home/heupr/golang/bin/core/pipeline/frontend
cp ./heupr.key /home/heupr/golang/bin/core/pipeline/frontend/cmd/heupr.key
cp ./heupr_io.crt /home/heupr/golang/bin/core/pipeline/frontend/cmd/heupr_io.crt
