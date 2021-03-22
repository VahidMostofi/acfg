## Environment variables to use during the process
```
export CLUSTER_NAME=bookstore-fargate
export AWS_REGION=us-west-2
export FARGATE_PROFILE_NAME=bookstore
export FARGATE_NAMESPACE=bookstore
export ACCOUNT_ID=$(aws sts get-caller-identity | cat | jq .'Account' -r)
export AWS_LOAD_BALANCER_CONTROLLER_NAME=aws-lb-controller
```
1. ### Create the cluster
    ```
    eksctl create cluster --name $CLUSTER_NAME --region $AWS_REGION --fargate --alb-ingress-access
    ```

2. ### Create Fargate Profile
    ```
    eksctl create fargateprofile --cluster $CLUSTER_NAME --name $FARGATE_PROFILE_NAME  --namespace $FARGATE_NAMESPACE
    ```

3. ### Add an IAM user to work with cluster
    ```
    eksctl utils associate-iam-oidc-provider \
        --region $AWS_REGION \
        --cluster $CLUSTER_NAME \
        --approve
    ```

4. ### Create IAM policy for load balancer
    ```
    wget https://raw.githubusercontent.com/kubernetes-sigs/aws-load-balancer-controller/main/docs/install/iam_policy.json

    aws iam create-policy --policy-name AWSLoadBalancerControllerIAMPolicy --policy-document file://iam_policy.json
    ```

5. ### Create IAM service account
    ```
    eksctl create iamserviceaccount \
    --cluster $CLUSTER_NAME \
    --namespace kube-system \
    --name $AWS_LOAD_BALANCER_CONTROLLER_NAME \
    --attach-policy-arn arn:aws:iam::${ACCOUNT_ID}:policy/AWSLoadBalancerControllerIAMPolicy \
    --override-existing-serviceaccounts \
    --approve
    ```

6. ### Apply loadbalancer crds
    ```
    run this somewhere else:
    git clone https://github.com/aws/eks-charts
    ```
    ```
    kubectl apply -k eks-charts/stable/aws-load-balancer-controller/crds
    ```

7. ### Make sure to have charts
    ```
    helm repo add eks https://aws.github.io/eks-charts
    ```

8. ### Check LBC VERSION
    ```
    if [ ! -x ${LBC_VERSION} ]
    then
        tput setaf 2; echo '${LBC_VERSION} has been set.'
    else
        tput setaf 1;echo '${LBC_VERSION} has NOT been set.'
    fi
    ```
    ```
    If the result is ${LBC_VERSION} has NOT been set., click here for the instructions.
    https://www.eksworkshop.com/020_prerequisites/k8stools/#set-the-aws-load-balancer-controller-version

    or just this: 
    export LBC_VERSION="v2.0.0"
    ```

9. ### Use helm to enable and upgrade aws load-balancer
    ```
    export VPC_ID=$(aws eks describe-cluster \
                    --name $CLUSTER_NAME \
                    --query "cluster.resourcesVpcConfig.vpcId" \
                    --output text)

    helm upgrade -i aws-load-balancer-controller \
        eks/aws-load-balancer-controller \
        -n kube-system \
        --set clusterName=${CLUSTER_NAME} \
        --set serviceAccount.create=false \
        --set serviceAccount.name=${AWS_LOAD_BALANCER_CONTROLLER_NAME} \
        --set image.tag="${LBC_VERSION}" \
        --set region=${AWS_REGION} \
        --set vpcId=${VPC_ID}
    ```
10. ### Deploy the bookstore application
    ```
    kubectl apply -f ns.yaml
    kubectl apply -f bookstore-components/
    ```

11. ### Deploy the endpoints-gateway and k8s-monitoring system
    - Change the template files in `monitoring-systems`.
    
    - The namespace in the k8s-monitor kube_inventory should be `bookstore`

    - Update the influxdb-secrets.yml
    ```
    kubectl apply -f monitoring-systems/metrics-server-components.yml
    kubectl apply -f monitoring-systems/influxdb-secrets.yml
    kubectl apply -f monitoring-systems/endpoint-gateway.yml
    kubectl apply -f monitoring-systems/k8s-monitor.yml
    ```

12. ### Create the ingress network
    ```
    kubectl apply -f ingress-network.yaml
    ```
13. ### Check the deployed ingress load balancer and the address

    ```
    kubectl get ingress/ingress-bookstore -n ${FARGATE_NAMESPACE}
    ```

14. ### Use the address field from step 13 to update the load generator.
    Don't forget `http://`

15. ### Update config file
    - Add `dev.env` to the root of the project. Use `dev.env.template` as template.
    - Update configs in `sample-configs/bookstore-aws.yml`

16.  ### Cleanup TODO
    ```
    kubectl delete -f monitoring-systems/metrics-server-components.yml
    kubectl delete -f monitoring-systems/influxdb-secrets.yml
    kubectl delete -f monitoring-systems/endpoint-gateway.yml
    kubectl delete -f monitoring-systems/k8s-monitor.yml
    kubectl delete -f bookstore-components/
    kubectl delete -f ns.yaml

    helm uninstall aws-load-balancer-controller -n kube-system

    eksctl delete iamserviceaccount \
        --cluster $CLUSTER_NAME \
        --name aws-load-balancer-controller \
        --namespace kube-system \
        --wait

    aws iam delete-policy --policy-arn arn:aws:iam::${ACCOUNT_ID}:policy/AWSLoadBalancerControllerIAMPolicy

    kubectl delete -k eks-charts/stable/aws-load-balancer-controller/crds

    eksctl delete fargateprofile \
    --name $FARGATE_PROFILE_NAME \
    --cluster $CLUSTER_NAME

    aws cloudformation delete-stack --stack-name eksctl-$CLUSTER_NAME-cluster
    ```