# IP Address Management (IPAM)

To avoid IP address conflicts and ensure efficient allocation of IP addresses, we need to keep track of ranges and their purpose.

## Address prefixes

| CIDR            | Purpose          | Description                                                         |
| --------------- | ---------------- | ------------------------------------------------------------------- |
| `172.31.0.0/16` | Cluster Pods     | Cluster pods for internal communication.                            |
| `172.30.0.0/16` | Cluster Services | Cluster services for internal service discovery and load balancing. |
| `172.29.0.0/16` | Load Balancers   | Externally accessible services of type `LoadBalancer`.              |
| `172.28.0.0/16` | Router IDs       | Router IDs as used by the routing protocol (BGP).                   |

## Static addresses

| IP Address    | Purpose  | Description                                  |
| ------------- | -------- | -------------------------------------------- |
| `172.30.0.10` | Core DNS | Core DNS server for cluster name resolution. |
| `172.28.0.1`  | alfa     | `lo` IP address of the `alfa` router.        |

## Autonomous System Numbers (ASN)

| ASN     | Purpose | Description                  |
| ------- | ------- | ---------------------------- |
| `65001` | `alfa`  | ASN for the `alfa` router.   |
| `65002` | `cph02` | ASN for the `cph02` cluster. |
