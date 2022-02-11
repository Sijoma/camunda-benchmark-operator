# Camunda-Benchmark-Operator 🏋️‍♀️

An operator that helps you perform benchmarks.  

## Your first benchmark
*This requires that you know how to run the operator, you might be able to figure that out from the <b>config</b> folder.*

In order to benchmark you need to create a secret with the credentials for your zeebe cluster.

An example secret is depicted below. Make sure to use the same keys in stringData.
If you change the name you also have to adjust the "credentialsSecretName" in the benchmark definition.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: cloud-credentials
type: Opaque
stringData:
  clientId: <your-client-id>
  clientSecret: <your-client-secret>
  zeebeAddress: <your-zeebe-address-with-port>
  authServer: <your-auth-server-url>
```

To now create a benchmark you apply the following manifest to your cluster. 
Make sure that the name of the secret matches with the secret you created earlier.

```yaml
apiVersion: cloud.camunda.io/v1alpha1
kind: Benchmark
metadata:
  name: benchmark-sample
spec:
  credentialsSecretName: cloud-credentials
  workerCount: 2
  processStarterRate: 200
  starterReplicas: 3
  duration: 1m30s
```


## Current Status

Run the following command if you want to get more insights on the current status of the benchmark:
`kubectl describe benchmark benchmark-sample`
```shell
Spec:
  Credentials Secret Name:  cloud-credentials
  Duration:                 2m
  Process Starter Rate:     200
  Starter Replicas:         3
  Worker Count:             2
Status:
  Progress:    Running
  Start Time:  2022-02-11T21:34:21Z

Events:        <none>
```

You can also run `kubectl get bench benchmark-sample` and it will print this:
```shell
NAME                    # WORKERS   STARTER RATE   DURATION   AGE    PROGRESS
benchmark-sample-gets   2           200            2m         106s   Running
```

In the future this should also log the results of the benchmark. 

