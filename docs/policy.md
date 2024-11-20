
# Policy
Access Control for all paths and secrets is managed in CryptKeeper using policies.
This document describes the structure and usage of policies in the CryptKeeper system. Policies are written in HCL (HashiCorp Configuration Language) and define permissions for accessing paths, secrets, and other resources.

## Policy Structure

Policies are structured into blocks that define access control rules for **paths** and **secrets**. Each block specifies:
1. **Permissions**: Allowed actions on the resource.
2. **Deny Permissions**: Actions explicitly denied on the resource.
3. **Users**: Users granted access to the resource.
4. **Certificates**: Certificates granted access.
5. **Apps**: Applications authorized for the resource.
6. **Groups**: User groups authorized for the resource.

---

### Sample Policy

```hcl
path "/kv-store" {
  permissions = []
  deny_permissions = ["create", "rotate", "update", "delete"]
  users       = ["vandana"]
  certificates = []
  apps = []
  groups = []
}

path "/kv-store" {
  permissions = ["list", "read", "create", "update", "delete", "rotate"]
  deny_permissions = []
  users       = []
  certificates = []
  apps = ["app-50fee5e4-6eb6-49be-81b3-ef478f44d259"]
  groups = ["admins"]
}

secret "/foo" {
  deny_permissions = ["read"]
  deny_users       = ["vandana"]
}

secret "/foo-1" {
  deny_permissions = ["read"]
  deny_groups = ["admins"]
}
```

---

## Path Policies

Path policies define access control rules for specific paths or namespaces.

### Fields

| Field            | Description                                                                 |
|-------------------|-----------------------------------------------------------------------------|
| `path`           | Specifies the resource path (e.g., `/kv-store`).                           |
| `permissions`    | A list of actions explicitly allowed on this path.                         |
| `deny_permissions` | A list of actions explicitly denied, overriding `permissions`.           |
| `users`          | A list of users with access to the path.                                   |
| `certificates`   | A list of certificates that can access the path.                          |
| `apps`           | A list of application IDs with access to the path.                        |
| `groups`         | A list of groups granted access to the path.                              |

### Example

#### Denying Specific Actions for a User
```hcl
path "/kv-store" {
  permissions = []
  deny_permissions = ["create", "rotate", "update", "delete"]
  users       = ["vandana"]
  certificates = []
  apps = []
  groups = []
}
```
- **Path**: `/kv-store`
- **Users**: Vandana has access but is denied the actions `create`, `rotate`, `update`, and `delete`.

#### Full Access for an Admin Group
```hcl
path "/kv-store" {
  permissions = ["list", "read", "create", "update", "delete", "rotate"]
  deny_permissions = []
  users       = []
  certificates = []
  apps = ["app-50fee5e4-6eb6-49be-81b3-ef478f44d259"]
  groups = ["admins"]
}
```
- **Path**: `/kv-store`
- **Permissions**: Full access granted (`list`, `read`, `create`, etc.).
- **Apps**: Access granted to a specific application ID.
- **Groups**: The `admins` group has access.

---

## Secret Policies

Secret policies define granular access control for specific secrets.

### Fields

| Field              | Description                                                             |
|---------------------|-------------------------------------------------------------------------|
| `secret`           | Specifies the secret name or path (e.g., `/foo`).                      |
| `deny_permissions` | A list of actions explicitly denied.                                   |
| `deny_users`       | Users explicitly denied access to the secret.                         |
| `deny_groups`      | Groups explicitly denied access to the secret.                        |

### Example

#### Denying Read Access for a User
```hcl
secret "/foo" {
  deny_permissions = ["read"]
  deny_users       = ["vandana"]
}
```
- **Secret**: `/foo`
- **Deny Permissions**: `read` access is denied for Vandana.

#### Denying Read Access for a Group
```hcl
secret "/foo-1" {
  deny_permissions = ["read"]
  deny_groups = ["admins"]
}
```
- **Secret**: `/foo-1`
- **Deny Permissions**: `read` access is denied for the `admins` group.

---

## Permissions List

### Supported Actions
- `list`: List available resources.
- `read`: Read the content of the resource.
- `create`: Create a new resource.
- `update`: Update an existing resource.
- `delete`: Delete a resource.
- `rotate`: Rotate a resource (e.g., keys, certificates).

### Deny Overrides
- Deny rules take precedence over allow rules.
- If an action is in `deny_permissions`, it cannot be performed even if listed in `permissions`.

---

## Policy Evaluation Rules

1. **Deny Takes Precedence**:
   - Actions listed in `deny_permissions` are blocked regardless of other settings.
2. **Explicit Permissions**:
   - Only actions listed in `permissions` are allowed.
3. **User and Group Scope**:
   - Access is determined by the intersection of `users`, `certificates`, `apps`, and `groups`.

---

## Best Practices

- Use **deny_permissions** for sensitive actions (e.g., `delete`, `rotate`).
- Leverage **groups** to simplify access management for large teams.
- Combine **path** and **secret** policies for layered security.

