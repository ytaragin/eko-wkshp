
kubectl create namespace workshop


cd tasks

docker build -t tasks-container:l2 .
kind load docker-image tasks-container:l2 --name workshop

cd ..
cd zvm-tunnel
docker build -t tunnel-container:l1 .
kind load docker-image tunnel-container:l1 --name workshop

cd ../protdb
docker build -t pg-migrate:v1 .     
kind load docker-image pg-migrate:v1 --name workshop



cd ..
helm upgrade --install wkshp protection-workshop --render-subchart-notes
