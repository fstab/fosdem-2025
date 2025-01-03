Zero-Code Distributed Traces for any programming language
---------------------------------------------------------

Demo for our FOSDEM 2025 presentation.

The following script will create a [kind](https://kind.sigs.k8s.io/) cluster and deploy the demo there.

```
./scripts/run.sh
```

Now, run the following command to expose port 3000 of the Grafana Pod on localhost:

```
kubectl port-forward $(kubectl get pods -lapp=grafana -o=name) 3000:3000
```

Access Grafana on [http://localhost:3000](http://localhost:3000). Default username is _admin_ with password _admin_.
