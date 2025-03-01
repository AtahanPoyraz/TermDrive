[![GPLv3 License](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://opensource.org/licenses/)

# TermDrive
TermDrive is a file storage and management server designed to be controlled from the terminal, providing a simple way for users to store and manage their files. It is designed to be lightweight and straightforward, offering a robust file sharing solution that supports various operations such as user registration, login, file upload, download, listing directory contents, and file deletion.

### Requirements
Before getting started, make sure you have the following tools installed:

```
* Go (Golang): Required for building and running the TermDrive server.
* Docker: For containerizing the application and making it easy to deploy.
* Make: A build automation tool used to simplify the setup process and manage the project.
```

### Key Features
#### **JWT-Based User Authentication**
TermDrive uses JSON Web Tokens (JWT) for secure user authentication. After a successful registration and login, users receive a token, which is used to authenticate subsequent requests and ensure that only authorized users can access their data. The token's default expiration time is 24 hours (which can be configured in the `config.yaml` file).

#### **File Storage and Management**
Users can upload, download, list directory contents, and delete files or directories on the server. All of these operations are secured by the JWT-based authentication system.

#### **Personal Directory**
Each user has their own personal directory within the `Storage` folder. This directory is automatically created when the user registers and provides a private space for their files.

#### **Command-Line Interface**
TermDrive is entirely controlled via the terminal using `curl` commands, making it ideal for users who prefer a terminal-based interface or need a lightweight file management solution.

#### **Secure and Flexible**
JWT usage provides secure access to the server, ensuring that users can only interact with their own directory, preventing unauthorized access. Additionally, the configuration file (`config.yaml`) allows for customization, offering flexibility in how the system is set up.

### How It Works
1. **User Registration**: Users register by providing a username, email, and password. The system automatically creates their personal directory in the `Storage` folder.

2. **Login and Token Generation**: After logging in, users receive a JWT token, which is required for all subsequent requests. The token is valid for 24 hours, after which the user must log in again to obtain a new token.

3. **File Operations**: With a valid JWT token, users can upload, download, list directory contents, and delete files or directories.

TermDrive is designed to run in a Docker environment, making it easy to set up and deploy. Configuration settings such as JWT expiration and storage paths can be customized in the `config.yaml` file located in the `server` directory.

Designed to be controlled via the terminal, TermDrive provides a secure, efficient, and flexible file management solution.

---
## Setup
#### **Clone the TermDrive repository locally**
To begin, clone the TermDrive repository to your local machine using the following commands:

```bash
git clone https://github.com/AtahanPoyraz/TermDrive.git
cd TermDrive
```

#### **Configure the server settings**
Modify the `TermDrive/server/config.yaml` file for any required settings, such as storage paths and JWT settings.

#### **Install the Necessary Dependencies and Prepare Docker**
To set up the environment and prepare Docker containers, follow the instructions below based on your operating system:

```bash
make run-setup
```

#### **Create Admin**
To create an admin user, set up the environment, and prepare Docker containers, use the following command based on your operating system:

```
make create-admin USERNAME=johndoe EMAIL=johndoe@example.com PASSWORD=example
```

#### **Run TermDrive**
To start the TermDrive server, use the following command:

```bash
make run-server
```

---
## Authentication
After setting up the project, you can follow these steps:

#### **Registration**

```bash
curl -X POST \
'http://localhost:8000/api/v1/auth/sign-up' \
-d '{"username":"johndoe", "email":"johndoe@example.com", "password":"example"}'
```

__Note__: This request creates a user account and sets up the user's personal directory within the `Storage` directory.

---
#### **Login**

```bash
curl -X POST \
'http://localhost:8000/api/v1/auth/sign-in' \
-d '{"email":"johndoe@example.com", "password":"example"}'
```

__Note__: This request allows the user to obtain a special token for further actions. The token will be in the format:

`Authorization: Bearer <JWT>`

With this token, you can continue your operations for 24 hours (this duration can be changed in `server/config.yaml`). After that, you will need to log in again to obtain a new token.

---
#### **Me**

```bash
curl -X GET \
'http://localhost:8000/api/v1/auth/me' \
-H "Authorization: Bearer <JWT>"
```

__Note__: This request retrieves the authenticated user's details from the server. It requires a valid JWT token for authorization.

---
## Usage
__Note__: All requests in this section require a valid authorization token. You must authenticate and obtain a token by logging in before performing any of the following operations.

#### **Upload File**

```bash
curl -X POST \
-F 'file=@/home/user/Desktop/path/your/file/example.txt' \
'http://localhost:8000/api/v1/storage/upload?path=path/your/file/example.txt' \
-H "Authorization: Bearer <JWT>"
```
__Explanation:__ This request uploads a file to your server. The `-F` flag specifies the file you want to upload (in this case, example.txt), and the path query parameter indicates the destination path where the file will be saved (path/your/file/example.txt). Make sure to replace the file path and destination path accordingly.

---
#### **Download File**

```bash
curl -X GET \
-f \
-o '/home/user/Desktop/path/your/file/example.txt' \
'http://localhost:8000/api/v1/storage/download?path=path/your/file/example.txt' \
-H "Authorization: Bearer <JWT>"
```

__Explanation:__ This request downloads a file from your server. The `-f` flag makes curl fail silently on errors, meaning if the download fails, no error message will be shown, and the target file won't be downloaded. It's useful in scripts to handle errors quietly. The `-o` flag specifies the local file path where the downloaded file will be saved (in this case, example.txt). The path query parameter indicates the file path to download from the server (path/your/file/example.txt).

---
#### **List Specific Path**

```bash
curl -X GET \
'http://localhost:8000/api/v1/storage/list?path=path/your/directory' \
-H "Authorization: Bearer <JWT>
```

__Explanation:__ This request lists the contents of a specific directory. The path query parameter indicates the directory you want to list. Replace path/your/directory with the directory path you wish to explore. The response will contain a list of files and folders within that directory.

---
#### **Delete Specific Path**

```bash
curl -X DELETE \
'http://localhost:8000/api/v1/storage/delete?path=path/your/directory' \
-H "Authorization: Bearer <JWT>
```

__Explanation:__ This request deletes a file or directory from your server. The path query parameter specifies the file or directory you want to delete. Replace path/your/directory with the appropriate path. Be cautious when using this request, as the deletion is permanent.

---
## Admin Endpoints

These endpoints require an admin role and allow administrators to perform operations such as managing users, controlling server settings, and monitoring server health. The following administrative functions are essential for a comprehensive system administration approach:

---
#### **Fetch Users By Pagination**

```bash
curl -X GET \
'http://localhost:8000/api/v1/admin/user/get?limit=10&offset=0' \
-H "Authorization: Bearer <JWT>
```

__Explanation:__ This request fetches users using pagination based on the limit (the number of users per request) and offset (where to start fetching users). This method is useful for querying users efficiently in smaller chunks, especially when you have a large number of users.
```
* limit=10: Returns 10 users per request.
* offset=0: Starts from the first user in the database.
```

---
#### **Fetch User By User ID**

```bash
curl -X GET \
'http://localhost:8000/api/v1/admin/user/get?userId=00000000-0000-0000-0000-000000000000' \
-H "Authorization: Bearer <JWT>
```

__Explanation:__ This request fetches a user based on their userId, which is a unique UUID for each user. The request will return data for the user whose ID matches the provided userId.
```
* userId: The UUID of the user. This parameter allows you to retrieve the data for a specific user.
```

---
#### **Fetch User By Username**

```bash
curl -X GET \
'http://localhost:8000/api/v1/admin/user/get?username=johndoe' \
-H "Authorization: Bearer <JWT>
```

__Explanation:__ This request fetches a user based on their username. The username parameter is used to find a specific user by their chosen username.
```
* username=johndoe: Retrieves the user with the username johndoe.
```

---
#### **Fetch User By Email**

```bash
curl -X GET \
'http://localhost:8000/api/v1/admin/user/get?email=johndoe@example.com' \
-H "Authorization: Bearer <JWT>
```

__Explanation:__ This request fetches a user based on their email. The email parameter is used to find a user by their registered email address.
```
* email=johndoe@example.com: Retrieves the user with the email address johndoe@example.com.
```

---
#### **Create User**

```bash
curl -X POST \
'http://localhost:8000/api/v1/admin/user/create' \
-d '{"username":"johndoe", "email":"johndoe@example.com", "password":"example", "role":"role", "isActive":true}' \
-H "Authorization: Bearer <JWT>"
```

__Explanation:__ This request allows the administrator to create a new user in the system. The provided data includes the following fields:
```
* username: The desired username for the new user.
* email: The email address for the new user.
* password: The password for the new user (make sure it meets the system's password requirements, e.g., complexity and length).
* role: The role assigned to the user. This can be "USER", "ADMIN".
* isActive: A boolean value (true or false) indicating whether the account should be active upon creation.
```

---
#### **Update User By User ID**

```bash
curl -X PATCH \
'http://localhost:8000/api/v1/admin/user/update?userId=00000000-0000-0000-0000-000000000000' \
-d '{"email":"johndoe@example.com", "password":"Passw0rq!1", "role":"role", "isActive":true}' \
-H "Authorization: Bearer <JWT>"
```

__Explanation:__ This request allows the administrator to update the details of an existing user in the system using their unique userId. The provided data includes the following fields:
```
* userId: UUID of the user to update.
* email: The new email address to be assigned to the user.
* password: The new password for the user (ensure it meets the system's password complexity requirements).
* role: The updated role of the user. It can be "USER", "ADMIN", etc.
* isActive: A boolean value (true or false) that indicates whether the user account should be active or deactivated.
```

---
#### **Delete User By User ID**

```bash
curl -X DELETE \
'http://localhost:8000/api/v1/admin/user/delete?userId=00000000-0000-0000-0000-000000000000' \
-H "Authorization: Bearer <JWT>
```
__Explanation:__ This endpoint allows the administrator to delete a user based on their unique user ID. When a user is deleted, all their associated data, including files and directories, will be removed. This operation is irreversible, so use with caution.
```
* userId: The UUID of the user. This parameter allows you to delete the data for a specific user.
```

---
#### **Fetch Files By Pagination**

```bash
curl -X GET \
'http://localhost:8000/api/v1/admin/storage/get?limit=10&offset=0' \
-H "Authorization: Bearer <JWT>
```

__Explanation:__ This request fetches files using pagination based on the limit (the number of files per request) and offset (where to start fetching files). This method is useful for querying files efficiently in smaller chunks, especially when you have a large number of files.
```
* limit=10: Returns 10 files per request.
* offset=0: Starts from the first user in the database.
```

---
#### **Fetch File By Filename And UserId**

```bash
curl -X GET \
'http://localhost:8000/api/v1/admin/storage/get?fileName=example&userId=00000000-0000-0000-0000-000000000000' \
-H "Authorization: Bearer <JWT>
```

__Explanation:__ This request fetches a file by file name and user ID. The file name and user ID parameters return the file for the registered user.
```
* fileName: The name of the file to retrieve from the server. 
* userId: The unique identifier of the user requesting the file. 
```

---
#### **Fetch File By Filename And UserId**

```bash
curl -X GET \
'http://localhost:8000/api/v1/admin/storage/get?filePath=example&userId=00000000-0000-0000-0000-000000000000' \
-H "Authorization: Bearer <JWT>
```

__Explanation:__ This request fetches a file by filename and user ID. The file name and user ID parameters are used to retrieve the file that belongs to the specified user.
```
* fileName: The name of the file to retrieve from the server.
* userId: The unique identifier of the user requesting the file.
```

---
#### **Fetch Files By User ID**

```bash
curl -X GET \
'http://localhost:8000/api/v1/admin/storage/get?userId=00000000-0000-0000-0000-000000000000' \
-H "Authorization: Bearer <JWT>
```

__Explanation:__ This request fetches a file by its file path and user ID. The file path and user ID parameters are used to retrieve the file associated with the specified user.
```
* userId: The UUID of the user. This parameter allows you to retrieve the data for a specific file.
```

---
#### **Fetch File By File ID**

```bash
curl -X GET \
'http://localhost:8000/api/v1/admin/storage/get?fileId=00000000-0000-0000-0000-000000000000' \
-H "Authorization: Bearer <JWT>
```

__Explanation:__ This request fetches a file based on their file ID, which is a unique UUID for each file. The request will return data for the file whose ID matches the provided file ID.
```
* fileId: The UUID of the file. This parameter allows you to retrieve the data for a specific file.
```

---
#### **Create File**

```bash
curl -X POST \
'http://localhost:8000/api/v1/admin/storage/create' \
-d '{"fileName":"example", "filePath":"example", "fileSize": 0, "mimeType":"example", "permissions": "-rwxr-xr-x","lastModifed": "exampleTimeStamp", "userId": "0000000-0000-0000-0000-000000000000"}' \
-H "Authorization: Bearer <JWT>"
```

__Explanation:__ This request allows the administrator to create a new file in the system. The provided data includes the following fields:
```
* fileName: The name of the file (string).
* filePath: The path where the file is stored (string).
* fileSize: The size of the file in bytes (integer).
* mimeType: The MIME type of the file (string).
* permissions: The file permissions (string) in the format of -rwxr-xr-x, representing read, write, and execute 
permissions for the file owner, group, and others.
* lastModified: Timestamp of the last modification time (ISO 8601 format, e.g., 2025-02-25T14:00:00Z).
* userId: UUID of the user who owns the file. This field identifies the user associated with the file.
```

---
#### **Update File By File ID**

```bash
curl -X PATCH \
'http://localhost:8000/api/v1/admin/storage/update?fileId=00000000-0000-0000-0000-000000000000' \
-d '{"fileName":"example", "filePath":"example", "fileSize": 0, "mimeType":"example", "permissions": "-rwxr-xr-x","lastModifed": "exampleTimeStamp", "userId": "0000000-0000-0000-0000-000000000000"}' \
-H "Authorization: Bearer <JWT>"
```

__Explanation:__ This request allows the administrator to update the details of an existing file in the system using their unique fileID. The provided data includes the following fields:
```
* fileId: UUID of the file to update. This is a required field and must correspond to an existing file in the system.
* fileName: The name of the file (string).
* filePath: The path where the file is stored (string).
* fileSize: The size of the file in bytes (integer).
* mimeType: The MIME type of the file (string).
* permissions: The file permissions (string) in the format of -rwxr-xr-x, representing read, write, and execute 
permissions for the file owner, group, and others.
* lastModified: Timestamp of the last modification time (ISO 8601 format, e.g., 2025-02-25T14:00:00Z).
* userId: UUID of the user who owns the file. This field identifies the user associated with the file.
```

---
#### **Delete File By File ID**

```bash
curl -X DELETE \
'http://localhost:8000/api/v1/admin/storage/delete?fileId=00000000-0000-0000-0000-000000000000' \
-H "Authorization: Bearer <JWT>
```
__Explanation:__ This endpoint allows the administrator to delete a file based on their unique file ID. 

```
* fileId: The UUID of the file. This parameter allows you to delete the data for a specific file.
```

---
## Feedback

If you have any feedback, please reach out to us at atahanpoyraz@gmail.com

---
## Authors

- [@AtahanPoyraz](https://www.github.com/AtahanPoyraz)
