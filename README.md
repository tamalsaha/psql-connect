$ kubectl create ns demo
namespace/demo created

$ kubectl create -f pg.yaml
postgres.kubedb.com/quick-postgres created

$ kubectl get pods -n demo
NAME               READY   STATUS    RESTARTS   AGE
quick-postgres-0   1/1     Running   0          43s


# use non-standard port to ensure it does not conflict with any locally running postgres
$ kubectl port-forward -n demo quick-postgres-0 5436:5432
Forwarding from 127.0.0.1:5436 -> 5432
Forwarding from [::1]:5436 -> 5432
Handling connection for 5436
Handling connection for 5436
Handling connection for 5436



## diff terminal window

$ kubectl get secret -n demo
NAME                                  TYPE                                  DATA   AGE
default-token-pc4w7                   kubernetes.io/service-account-token   3      93s
quick-postgres-auth                   Opaque                                2      91s
quick-postgres-snapshot-token-mcmpc   kubernetes.io/service-account-token   3      91s
quick-postgres-token-j2k64            kubernetes.io/service-account-token   3      91s

$ kubectl get secret -n demo quick-postgres-auth -o yaml
apiVersion: v1
data:
  POSTGRES_PASSWORD: LWFsVlVpakd2MzhMeHpTYQ==
  POSTGRES_USER: cG9zdGdyZXM=
kind: Secret
metadata:
  creationTimestamp: 2019-05-30T12:37:08Z
  labels:
    kubedb.com/kind: Postgres
    kubedb.com/name: quick-postgres
  name: quick-postgres-auth
  namespace: demo
  resourceVersion: "2090"
  selfLink: /api/v1/namespaces/demo/secrets/quick-postgres-auth
  uid: a3423d31-82d7-11e9-9768-0800271ce345
type: Opaque

$ echo -n 'LWFsVlVpakd2MzhMeHpTYQ==' | base64 -d
-alVUijGv38LxzSa

$ echo -n 'cG9zdGdyZXM=' | base64 -d
postgres


psql --help

Connection options:
  -h, --host=HOSTNAME      database server host or socket directory (default: "/var/run/postgresql")
  -p, --port=PORT          database server port (default: "5432")
  -U, --username=USERNAME  database user name (default: "tamal")
  -w, --no-password        never prompt for password
  -W, --password           force password prompt (should happen automatically)




https://stackoverflow.com/a/6405296

$ psql --host=127.0.0.1 --port=5436 --username=postgres
psql (11.3 (Ubuntu 11.3-1.pgdg18.04+1), server 10.2)
Type "help" for help.
postgres=#

# using docker
# https://stackoverflow.com/a/24326540

$ docker run --network=host -it  postgres:10.2 psql --host=127.0.0.1 --port=5436 --username=postgres
psql (10.2 (Debian 10.2-1.pgdg90+1))
Type "help" for help.

postgres=#





