# Amper - Web-Based Enterprise Software Platform
# Contact - info.grigor@gmail.com

**Project Status:** Showcase (Source Code Private) | **Development Period:** ~1 Year (Independent)

## Introduction

Amper is a comprehensive, web-based enterprise software solution developed over a year of dedicated independent work. It aims to provide businesses with an integrated suite of tools to manage their digital assets, communication, and custom data needs directly through a browser interface.

This repository serves as a detailed showcase of Amper's features and capabilities, reflecting the significant development effort undertaken.

# Video Showcase:
> 
[![IMAGE ALT TEXT](http://img.youtube.com/vi/qJ2LjMnXkBQ/0.jpg)](http://www.youtube.com/watch?v=qJ2LjMnXkBQ "Amper - Web-Based Enterprise Software Platform")

---

## Architecture & Scalability

Amper is designed with scalability at its core to handle growing data demands for organizations. The architecture separates the application logic from the data storage, enabling independent horizontal scaling:

* **Amper Nodes (Backend Services):** These nodes contain the core application logic, APIs, and business rules (implemented in Go). Multiple Amper nodes can be deployed and load-balanced to handle increasing user traffic and computational load.
* **Amper Datastore Nodes (Distributed Storage):** Large data objects (like Cloud Drive files, document versions, chat history) are stored in a separate layer of dedicated Datastore nodes. This layer functions as a distributed key-value store.
    * **Data Distribution:** When data needs to be stored, the Amper backend intelligently selects an available Datastore node. The key assigned to the data implicitly contains information about which Datastore node holds the actual value (the file/data blob).
    * **Unlimited Storage Scaling:** Administrators can dynamically add more Amper Datastore nodes to the system via a dedicated UI page within Amper (specifying IP address and port). This allows the total storage capacity to scale horizontally and virtually without limit, independent of the application nodes.
* **Scalability Management:** The Amper application includes administrative UI features for adding and managing both Amper (backend) and Amper Datastore nodes, providing control over the cluster's expansion.

This decoupled architecture ensures that both application performance and data storage capacity can be scaled independently based on specific organizational needs.

```mermaid
graph LR
    User[User Browser] --> LB(Load Balancer);
    LB --> A1["Amper Node 1 (Go)"]; 
    LB --> A2["Amper Node 2 (Go)"]; 
    LB --> A3["Amper Node ... (Go)"];

    subgraph Amper Backend Cluster
        A1; A2; A3;
    end

    A1 --> DB[(MySQL DB)];
    A2 --> DB;
    A3 --> DB;

    A1 --> DS1["Amper Datastore Node 1 (Go)"];
    A2 --> DS2["Amper Datastore Node 2 (Go)"];
    A3 --> DS3["Amper Datastore Node ... (Go)"];
    
    %% Optional: Show more links if desired, e.g.:
    %% A1 --> DS2; A2 --> DS3; A3 --> DS1; 

    subgraph Amper Datastore Cluster
        DS1; DS2; DS3;
    end

    style DB fill:#f9f,stroke:#333,stroke-width:2px
    style DS1 fill:#ccf,stroke:#333,stroke-width:2px,color:#000
    style DS2 fill:#ccf,stroke:#333,stroke-width:2px,color:#000
    style DS3 fill:#ccf,stroke:#333,stroke-width:2px,color:#000
```
---
## Core Features
    1. Advanced Cloud Drive
    2. Real-Time Communication Suite
    3. Dynamic Data Modeling & Management
    4. User Management
    5. Node management with horizontal scaling
---

## Technology Stack

Amper was built using the following technologies:

* **Frontend:**
    * Framework/Library: **React**
    * UI Components: **Material UI (MUI)**
    * Languages: JavaScript (ES6+), JSX, HTML5, CSS3
* **Backend:**
    * Language: **Go (Golang)**
* **Database:**
    * **MySQL** (Relational Database)
* **Communication Layer:**
    * UI Updates: **Asynchronous UI updates (simulating real-time) via AJAX polling** from the frontend to the backend.
* **Key Backend Integrations & Libraries:**
    * Document Processing: **LibreOffice** (invoked by the Go backend for Office document conversions)
    * Image Processing: **ImageMagick** (integrated via a Go library/wrapper)


## Core Features

Amper is built around several key modules designed to work together seamlessly:

### 1. Advanced Cloud Drive

Amper's Cloud Drive offers robust and intelligent file management capabilities beyond simple storage:

* **Hierarchical Structure:** Users can create nested folders to organize files logically.
* **File Operations:** Standard Cut, Copy, Paste, Delete, and Download functionalities are available.
* **Broad File Type Support:** Upload various file types, including documents, images, videos, and other common formats.
* **Intelligent Document Handling:**
    * **Automatic Conversion & Preview:** Uploaded Office documents (e.g., .docx, .xlsx, .pptx) are automatically converted to PDF.
    * **Thumbnail Generation:** A thumbnail preview (screenshot of the first page) is generated for quick identification.
    * **In-Browser PDF Rendering:** Documents are viewable directly within the browser using a built-in PDF viewer.
    * **Annotation & Editing:** Integration with an Adobe PDF editor plugin allows users to annotate and edit documents directly in the browser.
    * **Robust Version Control:** Every save/edit after annotation creates a new version of the document. All previous versions remain accessible, providing a complete history.
* **Image Processing:**
    * **Wide Format Compatibility:** Handles common image formats (PNG, JPEG, HEIC, TIFF, BMP, etc.).
    * **Web Optimization:** Images are converted to PNG format for reliable browser display.
    * **Thumbnail Generation:** Thumbnails are created for easy visual Browse in the file list.
* **Bulk Uploads:**
    * Users can drag-and-drop or select multiple files for simultaneous upload.
    * A dedicated sidebar panel shows the real-time progress of each file upload, including percentage completion.
* **File Up-Versioning:** Instead of uploading duplicates, users can download a file, modify it offline, and then upload it as a *new version* of the original file, maintaining a clear lineage and tracking progress effectively.

### 2. Real-Time Communication Suite

Amper includes a built-in chat system designed for instant and effective collaboration:

* **Direct Messaging (1-on-1):** Users can initiate private conversations with other individuals in the system.
* **Group Chat:** Users can create chat groups by inviting multiple participants for team or project discussions.
* **Channels:** Administrators can create large-scale channels and manage membership, suitable for announcements or organization-wide communication.
* **Live Updates (Real-Time Engine):**
    * **Instant Message Delivery:** New messages appear immediately for all participants without needing a browser refresh.
    * **Live Edits/Deletes:** Any modifications or deletions to messages are reflected instantly for users viewing the chat.
    * **Real-Time Reactions:** Users can react to messages with emojis, and these reactions (and their removal) appear live for everyone in the chat.
* **Persistent Chat History:** All conversations are saved and searchable.

### 3. Dynamic Data Modeling & Management

This is a core, powerful feature allowing businesses to tailor Amper to their specific data needs:

* **Admin-Defined Entities (Objects):** Administrators can define custom data objects relevant to their business (e.g., "Client," "Project," "Invoice").
* **Custom Fields:** For each object, administrators can define fields with specific data types:
    * Text
    * Number
    * Boolean (True/False)
    * Date
    * Date & Time
* **Object Types & Inheritance:**
    * Define different "types" for a base object (e.g., Object "Contact" could have Types "Lead," "Customer," "Partner").
    * Object Types can inherit fields from their parent Object Type, following Object-Oriented Programming (OOP) principles, allowing for structured specialization.
* **Database Integration:** Defined objects and their fields are automatically mapped to corresponding tables and columns in a MySQL database backend.
* **Dashboard & Widgets:** Users interact with the data through a configurable dashboard environment:
    * **Record Creation:** When creating a record for an object with defined types, the user is prompted to select the specific type, ensuring the correct set of fields is presented.
    * **Record List Widget:** Displays lists of records for a configured object. Features include:
        * CRUD Operations (Create, Read, Update, Delete).
        * Complex Filtering and Searching capabilities.
    * **Record Detail Widget:** Designed to show the full details of a single record.
        * **Widget Interaction:** Can be configured to "listen" to a Record List widget. Selecting a record in the List automatically displays its details in the linked Detail widget.
        * Allows users to conveniently view and update individual records.

### 4. User Management

Amper includes standard functionalities for managing user access:

* **User Registration:** Workflow for new users to sign up.
* **Email Activation:** System sends activation emails to verify user accounts.
* **Admin Management:** Administrators have tools to manage users (e.g., activate, deactivate, assign roles - *if roles were implemented*).
* **Secure Login:** Standard authentication process for users to access their accounts.

---

## Technology Stack Highlights

* **Backend Database:** MySQL
* **Architecture:** Web-Based (Accessible via Browser)
* *(Feel free to add Frontend/Backend languages/frameworks if you wish)*

---

## Project Status & Future Vision

Amper represents a significant solo development effort focused on building a versatile and integrated enterprise platform. While the source code is currently private, this repository serves as a comprehensive portfolio piece documenting its features.

Future development ideas included:

* **User Profile Pages:** Dedicated pages for users to manage their profile information.
* **Social Networking Elements:** Incorporating features like activity feeds, user walls, and connections to enhance collaboration and engagement within the platform.
* **Home Page:** A central landing page or dashboard summarizing key information and activities.

---

*This README provides a detailed overview of the Amper platform's features developed during its ~1-year incubation period.*
