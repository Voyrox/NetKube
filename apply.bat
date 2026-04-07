@echo off
setlocal

kubectl apply -f netkube.yaml
if errorlevel 1 goto :error

kubectl set image deployment/netkube netkube=glitchedking/netkube:latest -n netkube
if errorlevel 1 goto :error

kubectl rollout restart deployment/netkube -n netkube
if errorlevel 1 goto :error

kubectl rollout status deployment/netkube -n netkube
if errorlevel 1 goto :error

echo Deployment updated successfully.
exit /b 0

:error
echo Deployment update failed.
exit /b 1
