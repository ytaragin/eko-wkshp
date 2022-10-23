docker build -t prot-container:l1 .
kind load docker-image prot-container:l1 --name workshop
kubectl rollout restart deployment protection
