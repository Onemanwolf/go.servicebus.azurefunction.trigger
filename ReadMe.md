### Documentation for Configuring the Azure Function

This Azure Function is a **custom handler** that processes messages from an Azure Service Bus queue. Below is a detailed explanation of the app and how to configure it.

---

### **Overview of the App**
1. **Purpose**:
   - The app listens to messages from an Azure Service Bus queue (`queue`) and processes them.
   - If the message processing succeeds, it returns a `200 OK` response.
   - If the message processing fails (e.g., invalid message format or simulated failure), it returns a `500 Internal Server Error`, causing the message to be retried.

2. **Components**:
   - **`main.go`**: Contains the custom handler logic for processing Service Bus messages.
   - **`function.json`**: Defines the Service Bus trigger binding configuration.
   - **`host.json`**: Configures the Azure Functions runtime.

---

### **Code Explanation**

#### **`main.go`**
- **`queueHandler` Function**:
  - Decodes the incoming message from the Service Bus queue.
  - Processes the message and simulates success or failure.
  - Returns appropriate HTTP responses:
    - `200 OK` for successful processing.
    - `500 Internal Server Error` for failures (causing retries).

- **`main` Function**:
  - Starts an HTTP server to handle requests from the Azure Functions runtime.
  - The server listens on the port specified by the `FUNCTIONS_CUSTOMHANDLER_PORT` environment variable (default: `8080`).

---

#### **`function.json`**
Defines the Service Bus trigger binding:
```json
{
  "bindings": [
    {
      "name": "queueItem",
      "type": "serviceBusTrigger",
      "direction": "in",
      "queueName": "queue",
      "connection": "ServiceBusConnection",
      "autoComplete": false
    }
  ]
}
```
- **`type`**: Specifies the trigger type (`serviceBusTrigger`).
- **`queueName`**: The name of the Service Bus queue (`queue`).
- **`connection`**: The connection string name for the Service Bus (`ServiceBusConnection`).
- **`autoComplete`**: Set to `false`, meaning the function must explicitly handle message completion.

---

#### **`host.json`**
Configures the Azure Functions runtime:
```json
{
  "version": "2.0",
  "isDefaultHostConfig": true,
  "extensionBundle": {
    "id": "Microsoft.Azure.Functions.ExtensionBundle",
    "version": "[4.*, 5.0.0)"
  }
}
```
- **`version`**: Specifies the runtime version (`2.0`).
- **`extensionBundle`**: Ensures the required extensions for Service Bus are loaded.

---

### **Configuration Steps**

#### **1. Create an Azure Service Bus Queue**
1. Go to the Azure portal.
2. Create a Service Bus namespace.
3. Add a queue named `queue`.

#### **2. Configure the Azure Function**
1. **Set Up the Connection String**:
   - Add the Service Bus connection string to the local.settings.json file (for local development) or as an application setting in the Azure portal:
     ```json
     {
       "IsEncrypted": false,
       "Values": {
         "FUNCTIONS_WORKER_RUNTIME": "custom",
         "ServiceBusConnection": "<Your-Service-Bus-Connection-String>"
       }
     }
     ```

2. **Deploy the Function**:
   - Deploy the function to Azure using the Azure Functions Core Tools or your preferred deployment method.

3. **Set the `FUNCTIONS_CUSTOMHANDLER_PORT` Environment Variable**:
   - Ensure the `FUNCTIONS_CUSTOMHANDLER_PORT` environment variable is set to the port your custom handler listens on (default: `8080`).

#### **3. Configure the Service Bus Queue**
1. Set the **Max Delivery Count**:
   - Navigate to the Service Bus queue in the Azure portal.
   - Set the **Max Delivery Count** (default: `10`).
   - Messages exceeding this count will be moved to the dead-letter queue.

2. Set the **Lock Duration**:
   - Configure the lock duration to ensure the function has enough time to process messages.

---

### **Testing the Function**

#### **Local Testing**
1. Start the function locally:
   ```sh
   func start
   ```
2. Send a test message to the Service Bus queue using a tool like Azure Service Bus Explorer or the Azure SDK.

#### **Production Testing**
1. Deploy the function to Azure.
2. Send messages to the Service Bus queue and monitor the logs in the Azure portal.

---

### **Monitoring and Debugging**
1. **Logs**:
   - View logs in the Azure portal or locally using the Azure Functions Core Tools.
2. **Dead-Letter Queue**:
   - Inspect the dead-letter queue for messages that could not be processed.

---

### **Summary**
This Azure Function processes messages from a Service Bus queue using a custom handler. It is configured via function.json and host.json, and the message processing logic is implemented in main.go. Follow the steps above to configure, deploy, and test the function.