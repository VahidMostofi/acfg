1. ### Create the cluster
   _TODO clean up ports_
    ```
    k3d cluster delete cluster-k3d
    k3d cluster create cluster-k3d --servers 1 --agents 1 --port 9080:80@loadbalancer --port 9443:443@loadbalancer --port 9099:9099@loadbalancer --api-port 6443 --k3s-server-arg '--no-deploy=traefik'
    ```

2. ### Install & Configure ISTIO
    ```
    export PATH=$PATH:$HOME/.istioctl/bin
    istioctl install --set profile=default --set values.pilot.env.PILOT_HTTP10=1 -y
    kubectl label namespace default istio-injection=enabled
    ```
3. ### Deploy the bookstore application
    ```
    kubectl apply -f bookstore-components/
    ```

4. ### Deploy the endpoints-gateway and k8s-monitoring system
    - Change the template files in `monitoring-systems`.
    - Update the influxdb-secrets.yml
    ```
    kubectl apply -f monitoring-systems/metrics-server-components.yml
    kubectl apply -f monitoring-systems/influxdb-secrets.yml
    kubectl apply -f monitoring-systems/endpoint-gateway.yml
    kubectl apply -f monitoring-systems/k8s-monitor.yml
    ```

5. ### Create the ingress network
    ```
    kubectl apply -f ingress-network.yaml
    ```
6. ### Check the deployed ingress load balancer, and the address

    ```
    kubectl get service/istio-ingressgateway -n istio-system 
    ```

7. ### Use the address field from step 6 to update the load generator.
   _TODO_
   Don't forget `http://`

8. ### Update config file
   _TODO_
    - Add `dev.env` to the root of the project. Use `dev.env.template` as template.
    - Update configs in `sample-configs/bookstore-aws.yml`
