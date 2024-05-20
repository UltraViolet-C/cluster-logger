# cluster-logger

Take-home assignment made for Spectro Cloud's Jr. Software Developer interview process.

A basic kubernetes controller that reconciles a custom `ClusterScan` resource and logs the data from it in a nice, readable way. Mostly just an interface for encapsulating arbitrary jobs.

# How to test
I recommend running this using KIND to quickly run a Kubernetes cluster locally. The steps I followed are:

(all ran inside the project directory)

1. `kind create cluster` (if you are using `kind` to create a K8s cluster)
2. `make manifests`
3. `make install`
4. `make run`
5. (in a new terminal) `kubectl create -f config/samples/log_v1_clusterscan.yaml`
6. `kubectl get pods` -- Copy the name of the pod just created here.
7. `kubectl get logs [name of the pod from step 6]`

And voila! You will now see the logs generated from our sample `ClusterScan`.
