[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](../LICENSE)

# Digital Contracting Service

An automated orchestration workspace that deploys a [Digital Contracting Service](https://github.com/eclipse-xfsc/facis/tree/main/DCS) instance to a Kubernetes cluster.

---

## 🚀 Overview

The Digital Contracting Service (DCS) provides an open-source platform for creating, signing, and managing contracts digitally.
Integrated with the European Digital Identity Wallet (EUDI), it guarantees that all digital transactions are secure, legally binding, and interoperable.
DCS allows organizations to streamline business processes, reduce paperwork, and ensure compliance with eIDAS 2.0 regulations, while fostering trust across federated partners.

Key components of the Digital Contracting Service include:
- Multi-Contract Signing: Enables multi-party contract execution within a single integrated workflow.
- Automated Workflows: Automates contract generation, execution, and deployment to ensure legal
consistency and efficiency.
- Lifecycle Management: Monitors contracts with alerts for renewals, expirations, or required actions.
- Signature Management: Links contract signatures to verifiable digital identities to maintain legal validity
and trust.
- Secure Archiving: Stores signed contracts in a tamper-evident archive compliant with retention policies.
- Machine Signing: Supports automated signing for high-volume or routine transactions.

This module allows you to set up and interact with the Digital Contracting Service visually inside the ORCE environment. You don’t need to write any code or handle any complex API integration manually—just install the Node-RED node for Digital Contracting Service, drop it into your flow, and configure the endpoint and query.

Thanks to ORCE’s orchestration features, deploying a Digital Contracting Service instance and querying it happens in just a few clicks. Upload your configs, drag your node, and start the Digital Contracting Service

---

## ⚡️ Click-to-Deploy

---
## Prerequisites

Before running the deploy script, ensure you have the following:

### System Tools
The following CLI tools must be installed and accessible in your PATH:
- **kubectl** - Kubernetes command-line tool
- **helm** - Package manager for Kubernetes
- **jq** - Command-line JSON processor
- **curl** - Data transfer tool
- **sed** - Stream editor
- **openssl** - For generating TLS certificates and private keys
- **ssh-keygen** - For SSH key generation

### Kubernetes Cluster
- A working Kubernetes cluster
- **Traefik ingress controller** installed in the cluster (`kube-system` namespace)
  ```bash
  # Install Traefik (if not already installed)
  kubectl apply -f https://raw.githubusercontent.com/traefik/traefik-helm-chart/master/traefik/templates/deployment.yaml
  ```

### Files & Credentials
- **Kubeconfig file**: Path to your Kubernetes cluster configuration (e.g., `~/.kube/config`)
- **TLS Private Key**: Path to your TLS certificate private key (PEM format)
- **TLS Certificate**: Path to your TLS certificate (PEM format, must match your domain)
- **Domain**: A domain name where the DCS service will be accessible
- **URL Path**: A unique path identifier for this DCS instance

#### Generating TLS Certificates

If you don't have TLS certificates yet, you can generate self-signed ones using OpenSSL:

```bash
# Create a directory for certificates
mkdir -p ./certs

# Generate a private key
openssl genrsa -out ./certs/server.key 2048

# Generate a self-signed certificate (valid for 365 days)
openssl req -new -x509 -key ./certs/server.key -out ./certs/server.crt -days 365 \
  -subj "/CN=example.com/O=YourOrg/C=US"
```

**Note**: Replace `example.com` with your actual domain. For production, use certificates from a trusted Certificate Authority

### Keycloak Setup
- **Running Keycloak instance** with the following configured:
  - A configured realm
  - An `admin-cli` client with realm-admin role
  - A user with admin credentials for that realm
  - A client named `digital-contracting-service` with at least one role defined

#### Quick Keycloak Installation

```bash
# Deploy Keycloak
kubectl create -f https://raw.githubusercontent.com/keycloak/keycloak-quickstarts/refs/heads/main/kubernetes/keycloak.yaml

# Access via port-forward
kubectl port-forward svc/keycloak 8080:8080
```

Default credentials: `admin/admin` at http://localhost:8080

**Configure admin-cli client (Required):**
1. Login to Keycloak admin console
2. Select your realm (e.g., `dcs`)
3. Go to Clients → admin-cli
4. Go to Service account roles tab
5. Assign `realm-admin` role from `realm-management` client

---

## 🔐 Keycloak Configuration (Before Deployment)

Before running the deploy script, you must configure Keycloak with the required realm and client:

### 1. Create a Realm
1. Log in to Keycloak admin console (e.g., http://localhost:8080)
2. Click the realm dropdown (top-left)
3. Click "Create Realm"
4. Enter realm name (e.g., `dcs`)
5. Click "Create"

### 2. Create the OIDC Client
1. In your realm, go to **Clients** (left menu)
2. Click **Create client**
3. Enter client ID (e.g., `digital-contracting-service`)
4. For **Client type**, select **OpenID Connect**
5. Click **Next**
6. Enable: **Client authentication**, **Authorization**, **Standard flow enabled**
7. Click **Save**

### 3. Configure Redirect URIs
1. In your client settings, find **Valid redirect URIs**
2. Add: `https://<your-domain>/<path>/*`
   - Example: `https://example.com/dcs/*`
3. Click **Save**

---

## 🛠️ How to Use

### 1. Prepare the environment and prerequisites
You'll need:
1.1. A Kubernetes cluster to host the child instances
1.2. A local ORCE as the parent to host the initial developing environment

### 1.1. Kubernetes
"Orchestration Engine" node requires a working Kubernetes cluster with ingress installed on it. Initiate a K8s cluster and install nginx-ingress on it using this command.
```bash
export KUBECONFIG=`<YOUR KUBECONFIG PATH>`
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.12.3/deploy/static/provider/cloud/deploy.yaml
```
You can learn more by reading the [official documentation](https://kubernetes.github.io/ingress-nginx/deploy/)
After this step, you can proceed to step 1.2 (Installing a local ORCE)

### 1.2. Local ORCE
Install ORCE as described in the [ORCE page](https://github.com/eclipse-xfsc/orchestration-engine):
```bash
docker run -d --name xfsc-orce-instance -p 1880:1880 leanea/facis-xfsc-orce:1.0.16
```
Go to [http://localhost:1880](http://localhost:1880).

### Install Digital Contracting Service Node
Click on "New Node" in the sidebar.

![new button](./docImage/add-new-node.jpg?raw=true)

Upload `node-red-contrib-digital-contracting-service-0.0.1.tgz` from this repository and install. Refresh to activate the node.


### 2. Run the Deploy Script

Once all prerequisites are in place, you can deploy the Digital Contracting Service using the deploy script:

```bash
./deploy.sh \
  <kubeconfig> \
  <private_key_path> \
  <crt_path> \
  <domain> \
  <path> \
  <realm> \
  <oidc_client_id>
```

**Parameters:**
- `<kubeconfig>` - Path to your Kubernetes config file (e.g., `~/.kube/config`)
- `<private_key_path>` - Path to your TLS private key file
- `<crt_path>` - Path to your TLS certificate file
- `<domain>` - Domain name for the DCS instance (e.g., `example.com`)
- `<path>` - URL path identifier (e.g., `mydcs`, will create `https://example.com/mydcs`)
- `<realm>` - Keycloak realm name
- `<oidc_client_id>` - OIDC client ID (e.g., `digital-contracting-service`)

**Environment Variables (optional):**
- `DOCKER_REGISTRY` - Docker registry URL (e.g., `h6s71ks6.c1.de1.container-registry.ovh.net`)
- `DOCKER_REPO` - Docker repository namespace (e.g., `facis`)
- `OIDC_ISSUER_URL` - OIDC issuer URL for the backend (defaults to in-cluster URL: `http://keycloak.default.svc.cluster.local:8080/auth/realms/<realm>`)

**Example:**
```bash
./deploy.sh \
  ~/.kube/config \
  ./certs/server.key \
  ./certs/server.crt \
  example.com \
  mydcs \
  dcs \
  digital-contracting-service
```

The script will:
1. Verify all required tools are installed
2. Validate the kubeconfig file
3. Check for Traefik ingress controller
4. Wait for ingress External-IP assignment
5. Deploy the Helm chart
6. Create TLS secrets
7. Configure Keycloak authentication
8. Create user and assign roles
9. Wait for the service to be ready

---

### 3. Install your node
Click on the "Install" tab. Then on the upload icon. The node will be successfully installed.
![step two (flow)](./docImage/newstep.png?raw=true)


### 4. Create your flow
Drag in an Inject node, the **Digital Contracting Service** node, and a Debug node. Connect them:

![step three (flow)](./docImage/create-your-flow.png?raw=true)


### 5. Name your instance and configure the node
Double-click on the Digital Contracting Service node to open the edit dialog.
In this step, you must choose a **Digital Contracting Service Name**. This will become your instance’s unique identifier, so it must be:
- Unique (not used by any other instance)
- Free of special characters (letters and numbers only)
For example, if you name it `mydcs`, it will be used internally for instance referencing and must remain distinct.
![step four (flow)](./docImage/step2.png?raw=true)


### 6. Provide your kubeconfig file
In this tab, you need to provide the **kubeconfig** file of your target Kubernetes cluster.
This file allows the DCS node to access your Kubernetes environment and deploy the DCS instance correctly.
![step five (flow)](./docImage/step3.png?raw=true)


### 7. Provide domain address and TLS credentials
In this tab, you must enter the **domain address** where the DCS will be accessible. You’ll also need to upload your **TLS certificate** and **private key**.

The final accessible URL is formed by combining this domain with the DCS instance name you set earlier. For example:
- Instance Name: `mydcs`
- Domain: `example.com`
- Resulting URL: `example.com/mydcs`
Make sure your TLS credentials match the provided domain.
![step six (flow)](./docImage/step4.png?raw=true)


### 8. Information tab
After the service is successfully deployed, you can switch to the **Information** tab.
Here, the final URL of your deployed catalogue instance will be shown—ready to be copied and used for access or integration.
Click **Done** and then **Deploy**. Activate the Inject node.
![step eight (flow)](./docImage/step7.png?raw=true)
You should see JSON output in the Debug panel, showing catalogue entries.

---

## ⚙️ Configuration

Before running:

1. **DCS URL**  
   Set the URL of your DCS instance.

2. **Query Parameters**  
   Provide any filters or search strings in the node editor or in `msg.payload`.

3. **Authorization Token (optional)**  
   The DCS endpoints require auth headers (Bearer token).

---

## 📁 Directory Contents
```
.
├── node-red-contrib-digital-contracting-service-0.0.1.tgz
├── DigitalContractingService.html
├── DigitalContractingService.js
├── package.json
```

- **node-red-contrib-digital-contracting-service-0.0.1.tgz**  
  Installable node package.

- **DigitalContractingService.html**  
  Node-RED UI form.

- **DigitalContractingService.js**  
  Backend logic to send API requests and return results.

- **package.json**  
  Metadata and dependencies.

---

## 📦 Dependencies

```json
"node": ">=14.0.0",
"node-red": ">=3.0.0"
```

---

## 🔗 Links & References

- [Digital Contracting Service - XFSC](https://github.com/eclipse-xfsc/facis/tree/main/DCS)


---

## � Troubleshooting

### Image Pull Errors with Rancher Desktop

If you're using **Rancher Desktop** and encounter `ImagePullBackOff` errors when deploying with private registries:

1. **Root Cause**: Rancher Desktop uses containerd which is isolated from Docker's credential store. Even if Docker can pull an image, containerd may not have access to the registry credentials.

2. **Solution**: Manually import the image into Rancher Desktop:
   - Build or pull the image locally: `docker pull <registry>/<image>:tag`
   - Open **Rancher Desktop GUI** → Images
   - Click **Import** and select the image
   - The image will now be available to Kubernetes

---

## �📝 License

This project is licensed under the Apache License 2.0. See the [LICENSE](../LICENSE) file for details.
