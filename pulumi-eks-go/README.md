# pulumi-eks-go

## Create Resources
```
pulumi up
pulumi down
```

## Get Cluster Context
```
aws eks update-kubeconfig --region us-east-2 --name pulumi-eks-go-cluster-bcb8aae
```

## Refresh
```
pulumi refresh
```

## Export Stack to file

```
pulumi stack export > stack-status-$(date '+%Y-%m-%d_%H:%M:%S').json
```

## Import Stack File
```
pulumi stack import --file stack-status-2022-04-02_23:47:34.json
```

## Destroy specific resource

```
pulumi destroy --target urn:pulumi:pulumi-eks-go::pulumi-eks-go::aws:iam/role:Role::pulumi-eks-gonodegroup-iam-role --target-dependents
pulumi destroy --target urn:pulumi:pulumi-eks-go::pulumi-eks-go::aws:ec2/eip:Eip::pulumi-eks-go-eip1 --target-dependents
```

