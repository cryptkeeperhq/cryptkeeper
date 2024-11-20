import React, { useState, useEffect } from 'react';
import { ListGroupItem, Badge, Card, Table, Col, Container, ListGroup, Row } from 'react-bootstrap';
import { FaCertificate, FaDatabase, FaEnvelope, FaList } from 'react-icons/fa';


export const ApprovalRequestHelp = () => {
    return (

        <Card className='bg-transparent'>
            <Card.Header as="h2">What is Approval Workflow for Secrets?</Card.Header>
            <Card.Body>


                <p>The approval workflow for secrets in CryptKeeper is designed to add an additional layer of security and accountability when managing sensitive information. Here's a brief overview of why and when this workflow would be used:</p>

                <h5>Why Use an Approval Workflow?</h5>

                <ol>
                    <li><strong>Enhanced Security:</strong> By requiring approvals, you ensure that no single individual can unilaterally create, modify, or delete secrets without oversight.</li>
                    <li><strong>Accountability:</strong> The workflow provides a clear audit trail of who approved changes and when, which is essential for compliance and security audits.</li>
                    <li><strong>Minimizing Risk:</strong> It reduces the risk of accidental or malicious modifications to critical secrets by involving multiple stakeholders in the approval process.</li>
                </ol>

                <h5>When to Use an Approval Workflow?</h5>

                <ol>
                    <li><strong>High-Sensitivity Secrets:</strong> For secrets that are critical to the organization's security and operations, such as encryption keys, database credentials, and API keys.</li>
                    <li><strong>Regulatory Compliance:</strong> When managing secrets that are subject to regulatory requirements mandating multi-party approvals and audit trails.</li>
                    <li><strong>Operational Changes:</strong> During major updates or changes to the infrastructure that require temporary adjustments to secret access or configuration.</li>
                </ol>



            </Card.Body>
        </Card>
    )
}
export const PolicyManagementHelp = () => {
    const samplePolicy = `policy "example-policy" {
    description = "This is an example policy"

    path "secret/*" {
        capabilities = ["create", "read", "update", "delete", "list"]
    }

    path "secret/dev/*" {
        capabilities = ["read", "list"]
    }

    group "team1" {
        capabilities = ["read", "list"]
    }

    user "user1" {
        capabilities = ["create", "update"]
    }

    approle "role1" {
        capabilities = ["create", "read"]
    }
}`

    return (
        <Card className='bg-transparent'>
            <Card.Header>How do Policies work?</Card.Header>
            <Card.Body>
                <p className='lead'>In CryptKeeper, a Policy is a set of rules that defines the permissions for accessing and managing secrets within specific paths. Policies help in enforcing fine-grained access control, ensuring that only authorized users, groups, or applications can perform certain actions on the secrets.</p>


                <h2>Components of a Policy:</h2>
                <ul>
                    <li>Name and Description: Every policy has a unique name and a description to explain its purpose and scope.</li>
                    <li>Rules: Rules are the core part of a policy. They define what actions can be performed by whom. Each rule specifies the allowed actions (read, write, delete) and the entities (users, groups, app roles) that are granted those permissions.</li>
                </ul>

                <h2>Key Features:</h2>
                <ul>
                    <li>Fine-Grained Access Control: Policies provide detailed control over who can access and manage secrets. You can define different levels of permissions for different users, groups, or applications within the same path.</li>

                    <li>Resource-Based Policies: Policies are associated with specific paths, similar to resource-based policies in S3 buckets. This ensures that the access control rules are applied only to the intended secrets.</li>

                    <li>User and Group Permissions: Users can be part of multiple groups, each with its own set of policies. This allows for flexible and scalable access control management.</li>

                    <li>App Roles: In addition to user and group permissions, policies can also be applied to application roles, enabling secure access for services and applications.</li>

                    <li>Dynamic Policy Assignment: Policies can be dynamically assigned to users, groups, or app roles based on the current requirements. This helps in managing access control efficiently without making code changes.</li>
                </ul>

                A sample policy looks like this.

                <pre className='bg-dark text-white p-3 rounded-2 mt-2'>
                    {samplePolicy}
                </pre>

            </Card.Body>
        </Card>
    )
}

export const RoleManagementHelp = () => {
    return (
        <Card className='bg-transparent'>
            <Card.Header>How does App Roles work?</Card.Header>
            <Card.Body>
                <p className='lead'>App roles in CryptKeeper provide a secure way for machines and applications to authenticate and access secrets. App roles are automated identities designed specifically for non-human entities, enabling secure and efficient secret management for services and applications.</p>

                <h2>Components of App Roles:</h2>

                <ul>
                    <li>App Roles: App roles are identities assigned to applications and services. Each app role can have specific policies assigned to it, defining the permissions for accessing and managing secrets.</li>

                    <li>Policies: Policies for app roles work similarly to user and group policies. They define the actions that the app role can perform on specific paths, ensuring that applications have the necessary permissions to operate securely.</li>
                </ul>

                <h2>Key Features:</h2>
                <ul>
                    <li>Automated Authentication: App roles provide automated authentication for applications, ensuring that secrets are accessed securely and efficiently without manual intervention.</li>

                    <li>Granular Access Control: Policies assigned to app roles allow for granular control over what actions an application can perform on secrets. This ensures that applications have only the permissions they need to function.</li>

                    <li>Secure Secret Management: By using app roles, organizations can manage secrets for applications in a secure manner, reducing the risk of unauthorized access.</li>

                    <li>Dynamic Role Assignment: App roles can be dynamically created and assigned policies based on the application's requirements. This flexibility allows for efficient management of access control as applications evolve.</li>
                </ul>

            </Card.Body>
        </Card>
    )
}

export const UserManagementHelp = () => {
    return (
        <Card className='bg-transparent'>
            <Card.Header>How does Users and Groups work?</Card.Header>
            <Card.Body>
                <p className='lead'>In CryptKeeper, users and groups form the foundation for managing access control and permissions. Users can belong to multiple groups, and each group can have specific policies assigned to it. This hierarchical structure ensures flexible and scalable access management.</p>

                <h2>Components of User and Groups:</h2>
                <ul>
                    <li>Users: Users are individual entities that can access and manage secrets. Each user has a unique identity and can be part of one or more groups.</li>
                    <li>Groups: Groups are collections of users. Policies can be assigned to groups, allowing all members of the group to inherit the permissions defined by those policies. This simplifies the management of permissions for a large number of users.</li>
                </ul>
                <h2>Key Features:</h2>
                <ul>
                    <li>Scalable Permission Management: By organizing users into groups, CryptKeeper allows administrators to manage permissions for multiple users efficiently. Changes to group policies automatically apply to all group members.</li>

                    <li>Hierarchical Access Control: Users can inherit permissions from multiple groups, enabling a hierarchical approach to access control. This is useful for organizations with complex access requirements.</li>
                    <li>Dynamic User-Group Association: Users can be dynamically added or removed from groups as needed. This ensures that access control remains flexible and adapts to organizational changes.</li>
                    <li>Fine-Grained Permissions: Groups can be assigned different policies, allowing for fine-grained control over who can access and manage secrets within specific paths.</li>
                </ul>
            </Card.Body>
        </Card>

    )

}


export const TransitEncryptionHelp = () => {
    return (
        <Card className='bg-transparent'>
            <Card.Header>What is Transit Encryption?</Card.Header>
            <Card.Body>
                <p>
                    The Transit engine provides encryption and decryption as a service. It is designed to encrypt data in transit.
                    Users can encrypt and decrypt data without storing it. Users can encrypt and decrypt data without storing it. The engine supports various encryption algorithms.


                </p>
                <h2>Supported Key Types</h2>
                <Table responsive >
                    <thead>
                        <tr>
                            <th>Key Type</th>
                            <th>Description</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td>FPE-FF1</td>
                            <td>Format-preserving encryption</td>
                        </tr>
                        <tr>
                            <td>aes128-gcm96</td>
                            <td>AES-GCM with a 128-bit AES key and a 96-bit nonce; supports encryption, decryption, key derivation, and convergent encryption</td>
                        </tr>
                        <tr>
                            <td>aes256-gcm96</td>
                            <td>AES-GCM with a 256-bit AES key and a 96-bit nonce; supports encryption, decryption, key derivation, and convergent encryption (default)</td>
                        </tr>
                        <tr>
                            <td>chacha20-poly1305</td>
                            <td>ChaCha20-Poly1305 with a 256-bit key; supports encryption, decryption, key derivation, and convergent encryption</td>
                        </tr>
                        <tr>
                            <td>ed25519</td>
                            <td>Ed25519; supports signing, signature verification, and key derivation</td>
                        </tr>
                        <tr>
                            <td>ecdsa-p256</td>
                            <td>ECDSA using curve P-256; supports signing and signature verification</td>
                        </tr>
                        <tr>
                            <td>ecdsa-p384</td>
                            <td>ECDSA using curve P-384; supports signing and signature verification</td>
                        </tr>
                        <tr>
                            <td>ecdsa-p521</td>
                            <td>ECDSA using curve P-521; supports signing and signature verification</td>
                        </tr>
                        <tr>
                            <td>rsa-2048</td>
                            <td>2048-bit RSA key; supports encryption, decryption, signing, and signature verification</td>
                        </tr>
                        <tr>
                            <td>rsa-3072</td>
                            <td>3072-bit RSA key; supports encryption, decryption, signing, and signature verification</td>
                        </tr>
                        <tr>
                            <td>rsa-4096</td>
                            <td>4096-bit RSA key; supports encryption, decryption, signing, and signature verification</td>
                        </tr>
                        <tr>
                            <td>hmac</td>
                            <td>HMAC; supporting HMAC generation and verification</td>
                        </tr>
                    </tbody>
                </Table>
            </Card.Body>
        </Card>


    )
}
export const PathManagementHelp = () => {
    return (
        <Card className='bg-transparent'>
            <Card.Header >What is Path?</Card.Header>
            <Card.Body>
            <p className='lead'>In CryptKeeper, a Path is a unique namespace or directory within the system where secrets are stored and managed. Each path is associated with a specific engine type, which determines how the secrets within that path are handled. Paths help in organizing and segregating secrets, making it easier to manage access controls and apply different secret management policies based on the engine type.</p>

                        <div className="card mb-2 engine-card">
                            <div className="card-body">
                                <FaList />
                                <h2>KV (Key-Value)</h2>
                                <p>The KV engine stores secrets as key-value pairs. It's suitable for storing static secrets like API keys, passwords, and other configuration data. Users can create, read, update, and delete secrets.</p>
                            </div>
                        </div>
                        <div className="card mb-2 engine-card">
                            <div className="card-body">
                                <FaEnvelope />
                                <h2>Transit</h2>
                                <p>The Transit engine provides encryption and decryption as a service. It is designed to encrypt data in transit. Users can encrypt and decrypt data without storing it. The engine supports various encryption algorithms.</p>
                            </div>
                        </div>
                    
                        <div className="card mb-2 engine-card">
                            <div className="card-body">
                                <FaDatabase />
                                <h2>Database</h2>
                                <p>The Database engine manages database credentials. It dynamically generates database credentials with a configurable lease time. It helps in managing database access securely and efficiently.</p>
                            </div>
                        </div>
                        <div className="card engine-card">
                            <div className="card-body">
                                <FaCertificate />
                                <h2>PKI (Public Key Infrastructure)</h2>
                                <p>The PKI engine generates and manages X.509 certificates. It can be used to create a Certificate Authority (CA) or intermediate CAs, and issue certificates to users or applications. It's essential for managing secure communication channels.</p>
                            </div>
                        </div>
                    
            </Card.Body>
        </Card>
    )
}
export const CreateSecretHelp = () => {
    return (

        <Card className='bg-transparent'>
            <Card.Header>Learn More</Card.Header>
            <Card.Body>
                <p className='lead'>In CryptKeeper, a Path is a unique namespace or directory within the system where secrets are stored and managed. Each path is associated with a specific engine type, which determines how the secrets within that path are handled. Paths help in organizing and segregating secrets, making it easier to manage access controls and apply different secret management policies based on the engine type.</p>

                <Container className='p-0'>
                    <div className="row p-0 engine-overview">
                        <div className="col-lg-12 col-md-12 col-sm-12">
                            <FaList />
                            <h2>KV (Key-Value)</h2>
                            <p>The KV engine stores secrets as key-value pairs. It's suitable for storing static secrets like API keys, passwords, and other configuration data. Users can create, read, update, and delete secrets.</p>

                        </div><div className="col-lg-12 col-md-12 col-sm-12">

                            <FaEnvelope />
                            <h2>Transit</h2>
                            <p>The Transit engine provides encryption and decryption as a service. It is designed to encrypt data in transit. Users can encrypt and decrypt data without storing it. The engine supports various encryption algorithms.</p>

                        </div>
                        <div className="col-lg-12 col-md-12 col-sm-12">

                            <FaDatabase />
                            <h2>Database</h2>
                            <p>The Database engine manages database credentials. It dynamically generates database credentials with a configurable lease time. It helps in managing database access securely and efficiently.</p>

                        </div><div className="col-lg-12 col-md-12 col-sm-12">

                            <FaCertificate />
                            <h2>PKI (Public Key Infrastructure)</h2>
                            <p>The PKI engine generates and manages X.509 certificates. It can be used to create a Certificate Authority (CA) or intermediate CAs, and issue certificates to users or applications. It's essential for managing secure communication channels.</p>

                        </div>
                    </div>
                </Container>
            </Card.Body>
        </Card>

    );
};

// export { CreateSecretHelp };

