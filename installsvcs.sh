
kubectl create namespace workshop


cd tasks

docker build -t tasks-container:l2 .
kind load docker-image tasks-container:l2 --name workshop

cd ..
cd zvm-tunnel
docker build -t tunnel-container:l1 .
kind load docker-image tunnel-container:l1 --name workshop

cd ..
helm install wkshp protection-workshop