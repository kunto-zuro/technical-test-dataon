# **Technical Test - Node Hierarchy Management**

A RESTful API built with Golang to manage hierarchical data (tree structure) of nodes. The API allows creating, updating, deleting, and retrieving nodes in a tree structure. It is built using the **Echo framework**, **GORM** for database interaction, and **PostgreSQL** as the database.

---

## **Table of Contents**
- [Project Structure](#project-structure)
- [Database Schema](#database-schema)
- [API Endpoints](#api-endpoints)
- [Setup Instructions](#setup-instructions)
- [Testing the API](#testing-the-api)
- [Known Issues](#known-issues)
- [Author](#author)

---

## **Project Structure**

```
.
├── config
│   └── database.go         # Database configuration
├── dto
│   └── response.go         # DTO (Data Transfer Object) for API responses
├── models
│   └── node.go             # Database model for the `nodes` table
├── routes
│   └── routes.go           # API route definitions
├── service
│   ├── node_handlers.go    # HTTP handlers for API routes
│   ├── node_service.go     # Business logic layer for nodes
│   └── service.go          # Service interfaces and utilities
├── go.mod                  # Go module definition
├── main.go                 # Application entry point
└── README.md               # Documentation file
```

---

## **Database Schema**

### **Database Name**
`tree_db`

### **Table Name**
`nodes`

### **Table Definition**
| Column Name   | Data Type                  | Constraints                           |
|---------------|----------------------------|---------------------------------------|
| `id`          | BIGSERIAL                 | Primary Key, Auto-increment           |
| `code`        | TEXT                      | Unique, Not Null                      |
| `name`        | TEXT                      | Not Null                              |
| `parent_id`   | BIGINT                    | Nullable, References `nodes(id)`      |
| `created_at`  | TIMESTAMP WITH TIME ZONE  | Not Null, Default `CURRENT_TIMESTAMP` |
| `updated_at`  | TIMESTAMP WITH TIME ZONE  | Not Null, Default `CURRENT_TIMESTAMP` |

**SQL to Create Table:**
```sql
CREATE TABLE nodes (
    id BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    parent_id BIGINT REFERENCES nodes(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

---

## **API Endpoints**

### **Base URL**
```
http://localhost:8123
```

### **1. Get Entire Tree**
- **Method:** `GET`
- **Endpoint:** `/tree`
- **Description:** Retrieves the entire node hierarchy in a tree format.
- **Response:**
  ```json
  [
      {
          "id": 1,
          "code": "BOD00001",
          "name": "Board of Directors",
          "parent_id": null,
          "created_at": "2024-11-26T10:00:00Z",
          "updated_at": "2024-11-26T10:00:00Z",
          "list_division": [
              {
                  "id": 2,
                  "code": "DVS00001",
                  "name": "Information Technology",
                  "parent_id": 1,
                  "created_at": "2024-11-26T12:18:46.012766+07:00",
                  "updated_at": "2024-11-26T12:18:46.012766+07:00",
                  "list_division": []
              }
          ]
      }
  ]
  ```

---

### **2. Get Node by ID**
- **Method:** `GET`
- **Endpoint:** `/tree/:id`
- **Path Parameter:**
    - `id`: The ID of the node.
- **Description:** Retrieves a specific node and its child hierarchy.
- **Response:**
  ```json
  {
    "id": 2,
    "code": "DVS00001",
    "name": "Information Technology",
    "parent_id": 1,
    "created_at": "2024-11-26T12:18:46.012766+07:00",
    "updated_at": "2024-11-26T12:18:46.012766+07:00",
    "list_division": []
  }
  ```

---

### **3. Create Node**
- **Method:** `POST`
- **Endpoint:** `/tree`
- **Description:** Creates a new node under a specific parent.
- **Request Body:**
  ```json
  {
      "code": "NEW001",
      "name": "New Node",
      "parent_id": 1
  }
  ```
- **Response:**
  ```json
  {
      "id": 3,
      "code": "NEW001",
      "name": "New Node",
      "parent_id": 1,
      "created_at": "2024-11-26T12:00:00Z",
      "updated_at": "2024-11-26T12:00:00Z"
  }
  ```

---

### **4. Update Node**
- **Method:** `PUT`
- **Endpoint:** `/tree/:id`
- **Path Parameter:**
    - `id`: The ID of the node to update.
- **Description:** Updates a node's details, including its parent.
- **Request Body:**
  ```json
  {
      "code": "UPDATED001",
      "name": "Updated Node Name",
      "parent_id": 2
  }
  ```
- **Response:**
  ```json
  {
      "id": 3,
      "code": "UPDATED001",
      "name": "Updated Node Name",
      "parent_id": 2,
      "created_at": "2024-11-26T12:00:00Z",
      "updated_at": "2024-11-26T12:30:00Z"
  }
  ```

---

### **5. Delete Node**
- **Method:** `DELETE`
- **Endpoint:** `/tree/:id`
- **Path Parameter:**
    - `id`: The ID of the node to delete.
- **Description:** Deletes a node and all its children recursively.
- **Response:**
  ```json
  {
      "message": "Node and its children deleted successfully"
  }
  ```


---

### **6. Bulk Insert Nodes**
- **Method:** `POST`
- **Endpoint:** `/tree/bulk-insert`
- **Description:** Inserts multiple nodes in a single request while maintaining the parent-child hierarchy.
- **Request Body:**
  ```json
  {
    "id": 112,
    "code": "BOD00001",
    "name": "Board of Directors",
    "parent_id": 111,
    "created_at": "2024-11-26T12:18:46.012257+07:00",
    "updated_at": "2024-11-26T12:18:46.012257+07:00",
    "list_division": [
        {
            "id": 113,
            "code": "DVS00001",
            "name": "Information Technology",
            "parent_id": 112,
            "created_at": "2024-11-26T12:18:46.012766+07:00",
            "updated_at": "2024-11-26T12:18:46.012766+07:00",
            "list_division": [
                {
                    "id": 114,
                    "code": "ITSDERP",
                    "name": "ERP Development",
                    "parent_id": 113,
                    "created_at": "2024-11-26T12:18:46.012766+07:00",
                    "updated_at": "2024-11-26T12:18:46.012766+07:00",
                    "list_division": [
                        {
                            "id": 115,
                            "code": "ITSDERP-STAFF",
                            "name": "ERP Development STAFF",
                            "parent_id": 114,
                            "created_at": "2024-11-26T12:18:46.01398+07:00",
                            "updated_at": "2024-11-26T12:18:46.01398+07:00"
                        }
                    ]
                },
                {
                    "id": 116,
                    "code": "DVS00005",
                    "name": "Tech Development",
                    "parent_id": 113,
                    "created_at": "2024-11-26T12:18:46.01452+07:00",
                    "updated_at": "2024-11-26T12:18:46.01452+07:00",
                    "list_division": [
                        {
                            "id": 117,
                            "code": "DVS00005-STAFF",
                            "name": "Tech Development STAFF",
                            "parent_id": 116,
                            "created_at": "2024-11-26T12:18:46.015857+07:00",
                            "updated_at": "2024-11-26T12:18:46.015857+07:00"
                        }
                    ]
                },
                {
                    "id": 118,
                    "code": "DVS00008",
                    "name": "Software Maintenance",
                    "parent_id": 113,
                    "created_at": "2024-11-26T12:18:46.016364+07:00",
                    "updated_at": "2024-11-26T12:18:46.016364+07:00",
                    "list_division": [
                        {
                            "id": 119,
                            "code": "DVS00008-STAFF",
                            "name": "Software Maintenance STAFF",
                            "parent_id": 118,
                            "created_at": "2024-11-26T12:18:46.017202+07:00",
                            "updated_at": "2024-11-26T12:18:46.017202+07:00"
                        }
                    ]
                },
                {
                    "id": 120,
                    "code": "DVS00012",
                    "name": "Quality Assurance",
                    "parent_id": 113,
                    "created_at": "2024-11-26T12:18:46.017202+07:00",
                    "updated_at": "2024-11-26T12:18:46.017202+07:00",
                    "list_division": [
                        {
                            "id": 121,
                            "code": "DVS00012-STAFF",
                            "name": "Quality Assurance STAFF",
                            "parent_id": 120,
                            "created_at": "2024-11-26T12:18:46.017911+07:00",
                            "updated_at": "2024-11-26T12:18:46.017911+07:00"
                        }
                    ]
                },
                {
                    "id": 122,
                    "code": "DVS00007",
                    "name": "IT Support",
                    "parent_id": 113,
                    "created_at": "2024-11-26T12:18:46.018749+07:00",
                    "updated_at": "2024-11-26T12:18:46.018749+07:00",
                    "list_division": [
                        {
                            "id": 123,
                            "code": "DVS00007-STAFF",
                            "name": "IT Support STAFF",
                            "parent_id": 122,
                            "created_at": "2024-11-26T12:18:46.019433+07:00",
                            "updated_at": "2024-11-26T12:18:46.019433+07:00"
                        }
                    ]
                },
                {
                    "id": 124,
                    "code": "IMPOVS",
                    "name": "Implementation - AT",
                    "parent_id": 113,
                    "created_at": "2024-11-26T12:18:46.020047+07:00",
                    "updated_at": "2024-11-26T12:18:46.020047+07:00",
                    "list_division": [
                        {
                            "id": 125,
                            "code": "IMPOVS-STAFF",
                            "name": "Implementation - AT STAFF",
                            "parent_id": 124,
                            "created_at": "2024-11-26T12:18:46.020824+07:00",
                            "updated_at": "2024-11-26T12:18:46.020824+07:00"
                        }
                    ]
                },
                {
                    "id": 126,
                    "code": "DVS00014",
                    "name": "Human Resource",
                    "parent_id": 113,
                    "created_at": "2024-11-26T12:18:46.021935+07:00",
                    "updated_at": "2024-11-26T12:18:46.021935+07:00",
                    "list_division": [
                        {
                            "id": 127,
                            "code": "DVS00014-STAFF",
                            "name": "Human Resource STAFF",
                            "parent_id": 126,
                            "created_at": "2024-11-26T12:18:46.022537+07:00",
                            "updated_at": "2024-11-26T12:18:46.022537+07:00"
                        }
                    ]
                },
                {
                    "id": 128,
                    "code": "ITQA",
                    "name": "HR Development",
                    "parent_id": 113,
                    "created_at": "2024-11-26T12:18:46.023124+07:00",
                    "updated_at": "2024-11-26T12:18:46.023124+07:00",
                    "list_division": [
                        {
                            "id": 129,
                            "code": "ITQA-STAFF",
                            "name": "HR Development STAFF",
                            "parent_id": 128,
                            "created_at": "2024-11-26T12:18:46.023688+07:00",
                            "updated_at": "2024-11-26T12:18:46.023688+07:00"
                        }
                    ]
                },
                {
                    "id": 130,
                    "code": "DVS00013",
                    "name": "General Affairs",
                    "parent_id": 113,
                    "created_at": "2024-11-26T12:18:46.024369+07:00",
                    "updated_at": "2024-11-26T12:18:46.024369+07:00",
                    "list_division": [
                        {
                            "id": 131,
                            "code": "DVS00013-STAFF",
                            "name": "General Affairs STAFF",
                            "parent_id": 130,
                            "created_at": "2024-11-26T12:18:46.028104+07:00",
                            "updated_at": "2024-11-26T12:18:46.028104+07:00"
                        }
                    ]
                },
                {
                    "id": 132,
                    "code": "ITSFERP",
                    "name": "ERP Implementation",
                    "parent_id": 113,
                    "created_at": "2024-11-26T12:18:46.028611+07:00",
                    "updated_at": "2024-11-26T12:18:46.028611+07:00",
                    "list_division": [
                        {
                            "id": 133,
                            "code": "ITSFERP-STAFF",
                            "name": "ERP Implementation STAFF",
                            "parent_id": 132,
                            "created_at": "2024-11-26T12:18:46.029157+07:00",
                            "updated_at": "2024-11-26T12:18:46.029157+07:00"
                        }
                    ]
                },
                {
                    "id": 134,
                    "code": "DVS00016",
                    "name": "Training",
                    "parent_id": 113,
                    "created_at": "2024-11-26T12:18:46.029689+07:00",
                    "updated_at": "2024-11-26T12:18:46.029689+07:00",
                    "list_division": [
                        {
                            "id": 135,
                            "code": "DVS00016-STAFF",
                            "name": "Training STAFF",
                            "parent_id": 134,
                            "created_at": "2024-11-26T12:18:46.03022+07:00",
                            "updated_at": "2024-11-26T12:18:46.03022+07:00"
                        }
                    ]
                }
            ]
        },
        {
            "id": 136,
            "code": "DVS00003",
            "name": "Marketing and Sales",
            "parent_id": 112,
            "created_at": "2024-11-26T12:18:46.030904+07:00",
            "updated_at": "2024-11-26T12:18:46.030904+07:00",
            "list_division": [
                {
                    "id": 137,
                    "code": "DVS00017",
                    "name": "Sales",
                    "parent_id": 136,
                    "created_at": "2024-11-26T12:18:46.030904+07:00",
                    "updated_at": "2024-11-26T12:18:46.030904+07:00"
                }
            ]
        },
        {
            "id": 138,
            "code": "PURC",
            "name": "Purchasing",
            "parent_id": 112,
            "created_at": "2024-11-26T12:18:46.031466+07:00",
            "updated_at": "2024-11-26T12:18:46.031466+07:00"
        },
        {
            "id": 139,
            "code": "DVS00015",
            "name": "Finance",
            "parent_id": 112,
            "created_at": "2024-11-26T12:18:46.032495+07:00",
            "updated_at": "2024-11-26T12:18:46.032495+07:00"
        },
        {
            "id": 140,
            "code": "SPRJT0001",
            "name": "Special Project",
            "parent_id": 112,
            "created_at": "2024-11-26T12:18:46.033221+07:00",
            "updated_at": "2024-11-26T12:18:46.033221+07:00",
            "list_division": [
                {
                    "id": 141,
                    "code": "PMJ0001",
                    "name": "PM Jakarta",
                    "parent_id": 140,
                    "created_at": "2024-11-26T12:18:46.033727+07:00",
                    "updated_at": "2024-11-26T12:18:46.033727+07:00"
                },
                {
                    "id": 142,
                    "code": "DKR0001",
                    "name": "Dekorey",
                    "parent_id": 140,
                    "created_at": "2024-11-26T12:18:46.034385+07:00",
                    "updated_at": "2024-11-26T12:18:46.034385+07:00"
                }
            ]
        }
    ]
  }
  
  ```
- **Response:**
  ```json
  {
      "message": "Nodes inserted successfully"
  }
  ```

---

## **Setup Instructions**

### **1. Clone the Repository**
```bash
git clone https://github.com/kunto-zuro/technical-test-dataon.git
cd technical-test-dataon
```

### **2. Configure the Database**
Ensure PostgreSQL is installed and running. Create the database and table using the provided schema above. 
- Database name: tree_db
- Table name: nodes

### **3. Install Dependencies**
Install the required Go modules:
```bash
go mod tidy
```

### **4. Run the Server**
Start the server:
```bash
go run main.go
```

The server will start at `http://localhost:8123`.

---

## **Testing the API**

### **Using curl**
- **Get Tree**:
  ```bash
  curl --location --request GET 'http://localhost:8123/tree'
  ```

- **Create Node**:
  ```bash
  curl --location --request POST 'http://localhost:8123/tree' \
  --header 'Content-Type: application/json' \
  --data-raw '{
      "code": "NEW001",
      "name": "New Node",
      "parent_id": 1
  }'
  ```

---

### Penjelasan
Ketentuan: 
- Level sudah di atur dibagian API insert, insert bulk, & update. Tidak dapat melakukan insert kembali jika levelingnya sudah lebih dri 5 level.
- Code tidak boleh duplicate 
- Karena data list division sudah bertipe array sehingga JSON datanya dapat di expand ataupun di collapse
- Untuk level ke-5 dari setiap data DVS00001-Information Technology sudah saya buat dengan asumsi sebagai STAFF disetiap Divisi berada pada API bulk insert.
