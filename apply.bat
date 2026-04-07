@echo off
setlocal

kubectl set image deployment/netkube netkube=glitchedking/netkube:latest -n netkube
kubectl rollout restart deployment/netkube -n netkube
kubectl rollout status deployment/netkube -n netkube

endlocal