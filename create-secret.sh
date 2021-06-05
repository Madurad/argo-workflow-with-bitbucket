NAMESPACE="argo" # Namespace of argo-workflow

kubectl create secret -n $NAMESPACE generic bitbucket-creds --from-file=id_rsa=$1