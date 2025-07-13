# JavaScript Test Framework Improvements

Why: To provide better organization and more comprehensive testing capabilities

- Refactored main.go to separate step tests into a dedicated function
- Added conversation testing to validate the JS conversation wrapper
- Added command-line flags to selectively run different test suites
- Improved test output with clear section markers
- Added comprehensive conversation API testing

# Restructure RAG-to-Riches Application

Reorganized the application structure for better maintainability and separation of concerns.

- Created proper directory structure with internal packages
- Separated handlers into individual files
- Added mock document list view
- Improved template organization with subdirectories
- Enhanced CSS styling for consistency

# Schema Management API Implementation

Added HTMX-based API endpoints for schema management to support interactive schema creation, searching, and viewing.

- Added /api/v1/schemas/search endpoint for filtering schemas
- Added /api/v1/schemas/new endpoint for creating new schemas
- Added /api/v1/schemas/{id}/details endpoint for viewing schema details
- Implemented zerolog logging for all schema operations
- Added mock data store for development testing

# Schema Management Code Organization

Improved code organization and template management for schema functionality.

- Split HTMX templates into separate file for better maintainability
- Removed unused imports and fixed linter errors
- Improved template error handling and logging
- Centralized template definitions for schema-related views

# Schema Management UI Improvements

Enhanced schema management interface with HTMX integration.

- Added proper form handling for schema creation
- Implemented loading indicators for all async operations
- Added toast notifications for success/error messages
- Improved schema editor layout and functionality
- Added new schema form template with validation

# Schema Index Management

Added index creation and status monitoring for schemas.

- Implemented create-index endpoint with mock functionality
- Added index status display to schema details view
- Enhanced schema details UI with status cards
- Added loading indicators for index creation
- Implemented index health monitoring display

# Schema Management Tab Navigation

Improved schema details view with proper tab navigation.

- Added automatic tab switching when viewing schema details
- Implemented proper Bootstrap tab initialization
- Added event listeners for HTMX content loading
- Improved tab content organization and structure
- Enhanced user experience with smoother tab transitions

# Schema Management Implementation

Added core schema management functionality to support RAG document processing.

- Created schema package with CRUD operations
- Implemented file-based storage using YAML
- Added search functionality with filtering
- Added proper error handling and logging
- Added schema ID generation and update functionality
- Added GenerateID function using crypto/rand with timestamp fallback
- Added UpdateSchema method to Manager for handling mutable field updates
- Ensured schema definition remains immutable during updates

# Schema Management Integration

Integrated schema package with HTTP handlers and removed mock implementation.

- Replaced mock schema storage with file-based implementation
- Updated handlers to use schema manager
- Added schema manager to server struct
- Kept index management as mock for now pending ES integration

# Template Organization

Extracted HTMX templates from schemas_htmx.go into individual template files for better organization and maintainability. Templates are now located in templates/schemas/ directory.

- Moved schema list template to schema-list.html
- Moved schema details template to schema-details.html
- Moved success message template to success.html
- Moved new form template to new-form.html
- Moved index status template to index-status.html

# Schema Management Improvements

Enhanced the schema management functionality to fix issues with ID handling and updates.

- Fixed schema ID handling in list and detail views
- Added proper schema update functionality for mutable fields (name and notes)
- Added automatic list refresh after schema creation
- Improved error handling and success messages
- Made schema definition immutable after creation

# Schema List Refresh Fix

Why: To properly update the schema list after saving while maintaining selection

- Added selected state handling to schema list template
- Updated schema search handler to handle selected schema ID
- Modified list refresh to maintain current schema selection
- Added visual feedback for selected schema in list

# Remove Version field from Schema struct

Simplified the Schema struct by removing the Version field since it was causing issues and wasn't being used properly.

- Removed Version field from Schema struct
- Updated GetSchema to not require version parameter
- Updated schema handlers to work without version information

# Elasticsearch Index Creation Implementation

Added actual Elasticsearch index creation functionality to replace mock implementation.

- Created IndexCreationJob to handle index creation process
- Added proper error handling and logging for ES operations
- Integrated ES client with schema management system
- Used schema definition as index mapping
- Added default index settings with configurable shards/replicas

# Add Elasticsearch Client Configuration

Added Elasticsearch client initialization with environment variable configuration.

- Added ES client to Server struct
- Added environment variable configuration for ES connection (ELASTICSEARCH_URL, ELASTICSEARCH_USER, ELASTICSEARCH_PASSWORD)
- Added connection test during server initialization
- Added fallback to default localhost URL if not configured

# Real-time Elasticsearch Index Status

Added real-time Elasticsearch index status checking to schema details view.

- Added ES index stats retrieval for document count and size
- Added ES health check for index status
- Added proper size formatting (B, KB, MB, GB)
- Added nil safety for non-existent indices
- Added proper error handling for ES API calls

# Refactor Elasticsearch Helpers

Moved Elasticsearch functionality into dedicated helpers package for better reusability.

- Created helpers/es.go for shared ES functionality
- Moved index creation logic to reusable IndexCreationJob
- Moved index status checking to GetIndexStatus helper
- Improved error handling with more detailed error messages
- Added better response body handling for ES errors

# Navigation Improvements

Added active state highlighting to the main navigation bar to improve user experience and site navigation.

- Updated base.html template to conditionally show active state for current section
- Removed hardcoded active states
- Added dynamic active state based on template data

# Rename Search to Queries for Consistency

Renamed search-related functions and types to query-related ones for better consistency in terminology. This change makes the code more aligned with the domain language where we perform queries against the document store.

- Renamed SearchResult to QueryResult
- Renamed SearchPrompt to QueryPrompt
- Renamed search-related handler functions to query-related ones
- Updated route registration to use "queries" instead of "search"

# URL Updates for Schema Navigation

Why: To enable direct navigation to specific schemas via URL and improve user experience

- Added new route /schemas/{id} for direct schema access
- Updated schema list to modify URL when selecting schemas
- Added automatic schema details loading when accessing direct URLs

# Schema URL Navigation Fix

Why: To fix issues with schema URL navigation and improve the user experience

- Fixed schema ID handling in URL updates
- Improved direct schema URL navigation
- Updated schema list items to properly target details tab
- Simplified schema selection handling with single HTMX request

# Improved Schema URL Navigation

Why: To better utilize HTMX's built-in URL handling capabilities and improve link behavior

- Switched to hx-push-url for automatic URL updates
- Made schema links real URLs that work with browser navigation
- Simplified JavaScript code by removing manual URL manipulation
- Improved tab switching behavior with event listeners

# Schema List Template Consolidation

Why: To simplify template organization and fix URL handling

- Consolidated schema list items into main list.html template
- Removed unused schema-list.html template
- Updated schema links to use proper URLs with hx-push-url

# Direct Schema URL Navigation

Why: To ensure consistent rendering when accessing schema URLs directly

- Added HTMX request detection in HandleSchemaByID
- Added schema existence verification before rendering
- Unified template rendering between direct and HTMX requests
- Added proper 404 handling for non-existent schemas

# Schema Handler Consolidation

Why: To simplify code and improve handler organization

- Merged HandleSchemaByID and handleSchemaDetails into a single handler
- Moved schema details page route to registerSchemaRoutes
- Improved code organization and reduced duplication
- Maintained consistent behavior for both direct and HTMX requests

# Schema Creation URL Handling

Why: To ensure proper URL updates when creating new schemas

- Updated new schema form to push correct URL after creation
- Modified schema creation handler to return schema details directly
- Added HX-Push-Url header for proper URL updates
- Improved user experience with immediate schema details display

# Schema Creation Form Fix

Why: To correctly handle URL updates for new schemas

- Removed client-side hx-push-url from new schema form
- Rely on server-side HX-Push-Url header for URL updates
- Simplified form template

# Schema URL Convention Update

Why: To follow HTMX and REST conventions more closely

- Changed from path-based to query-based schema selection
- Simplified URL structure to use ?schema= parameter
- Removed /schemas/{id} route in favor of query params
- Improved URL sharing and bookmarking capability

# Schema List URL Handling

Why: To properly handle URL updates while maintaining API separation

- Updated schema list to use hx-push-url with query parameters
- Separated URL state from API endpoint fetching
- Updated schema detection to use URLSearchParams
- Maintained clean URLs while fetching from API endpoints

# Schema View Simplification

Why: To simplify the UI and improve usability

- Removed tab-based navigation in schema view
- Consolidated schema details into main view
- Simplified template structure and HTMX interactions
- Improved layout for better readability

# Schema Form Target Fix

Why: To fix schema saving functionality after view simplification

- Updated new schema form to target schema-details container
- Aligned form behavior with simplified view structure
- Removed unused message container references

# Schema Save Functionality Fix

Why: To fix non-working save button and improve error handling

- Added form-level HTMX indicator
- Improved error handling and validation in new schema handler
- Added debug logging for form submission
- Simplified save button structure

# Schema Search Endpoint Fix

Why: To fix schema list not updating after save operations

- Added fallback to ListSchemas when no search parameters provided
- Improved search handler logging and error handling
- Ensured consistent template usage for schema list
- Fixed schema list refresh functionality

# Schema Management UI Migration to Templ

Updated schema management components to use existing types and structures.

- Updated SchemaDetailsView to match existing schema types
- Fixed component naming and structure
- Improved type safety with proper imports
- Maintained existing functionality while migrating to templ

# Convert HTML templates to templ components

Converted remaining HTML templates to templ components for better type safety and consistency. Moved success template to components package for reuse.

- Converted success.html to templ component in components package
- Converted new-form.html to templ component
- Converted index-status.html to templ component
- Updated schemas.go to use new templ components
- Removed old HTML template references from schemas_htmx.go

# Convert Error Template to Templ Component

Converted error template to templ component for consistency and type safety.

- Created Error component in components package
- Updated all error template usages to use new templ component
- Removed schemas_htmx.go as all templates are now in templ

# Improve Base Template Parameters

Added a PageData struct to make base template parameters more explicit and maintainable.

- Created PageData struct with Title and Active fields
- Updated Base template to use PageData struct
- Added proper page titles to schema views

# Document Processing Prompt Management

Added proper prompt management for document processing prompts, replacing mock implementations.

- Integrated existing prompt manager with document handlers
- Added prompt type support (document-processing, query-transformation, result-synthesis)
- Maintained mock data for test cases and documents while implementing real prompt storage
- Updated handlers to use proper prompt storage and retrieval
- Added proper error handling and logging for prompt operations

# Prompt Manager Refactoring

Simplified prompt management by creating dedicated managers for each prompt type.

- Created separate prompt managers for document processing, query transformation, and synthesis
- Removed prompt type parameter from manager methods
- Moved system prompt to metadata field for flexibility
- Updated handlers to use dedicated document prompt manager
- Improved code organization with clearer separation of concerns

# Prompt Manager Type Enforcement

Improved prompt type safety by enforcing types at the manager level.

- Added Type field to Manager struct to store prompt type
- Removed type parameter from file paths since it's now in the base path
- Ensure prompt type matches manager type when reading files
- Simplified prompt creation by automatically setting type
- Removed redundant type checks and parameters

# Document Handler Schema Integration

Integrated real schema management into document handlers.

- Updated new prompt form to use real schemas from schema manager
- Added schema name lookup when displaying prompts
- Improved error handling for schema lookups
- Maintained mock data for test cases and documents
- Added proper logging for schema-related operations

# Document Prompt Form Improvements

Enhanced document prompt form handling and validation.

- Added proper template field handling from form submission
- Added validation for required fields (name, schema, template)
- Added debug logging for form values
- Improved error messages for missing fields
- Added monospace font for template textarea

# Document Prompt Update Functionality

Added proper update functionality for document prompts.

- Added form validation for prompt updates
- Added proper template and system prompt handling
- Added debug logging for update operations
- Improved error messages and validation
- Updated prompt details view with edit form
- Made schema field read-only in edit mode
- Added proper HTMX integration for updates

# Query Prompt Management UI

Added query prompt management UI with HTMX integration for managing and testing query transformation prompts. The UI allows creating, editing and testing prompts that transform natural language queries into structured search queries.

- Added query prompt list view with search and filtering
- Added query prompt creation and editing forms
- Added query prompt testing interface
- Added test history tracking and viewing
- Integrated with HTMX for dynamic updates

# Unified Test Results Display

Extended query test results to match document test results for consistent monitoring and debugging experience.

- Added execution time, creation time, success status, and token usage details to query TestResults type
- Updated query test results template to show comprehensive execution metrics
- Maintained consistent UI between document and query test results displays

# Query Test Results Consistency

Updated query test handlers to include all execution metrics for better consistency with document test results.

- Added execution time, success status, and token details to query test results
- Updated test case retrieval to include all execution metrics
- Maintained consistent test result structure between documents and queries

# Add Schema Support to Query Prompts

Added schema ID and name fields to query prompts to maintain consistency with document prompts and enable proper data validation against schemas.

- Added SchemaID and SchemaName fields to query.PromptData
- Updated query prompt forms to include schema selection
- Added schema validation in query prompt creation
- Modified query handlers to handle schema information

# Document Manager Integration and Multi-Document Testing

Added integration with document manager to show real documents in the test interface. Updated test interface to allow selecting multiple documents for testing prompts against multiple inputs at once.

- Integrated document manager to list real documents
- Updated test form to use checkboxes for multi-document selection
- Modified test results to show outcomes for multiple documents
- Added error handling for document retrieval

# Single Test Execution for Multiple Documents

Updated test execution to process multiple documents in a single test run for better efficiency and consistency.

- Modified TestResults struct to track multiple input documents
- Updated test execution to process all selected documents in one run
- Enhanced test results display to show per-document extraction results
- Added document name display in extraction results

# Document Handler Refactoring

Extracted document handling logic into a separate DocumentsHandler struct to improve code organization and maintainability.

- Created new document/handlers.go file to contain all document-related HTTP handlers
- Moved document handling logic from Server struct to DocumentsHandler
- Updated Server to use the new DocumentsHandler

# Document Processing Runner Extraction

Extracted document processing logic into a dedicated Runner type to improve separation of concerns and prepare for actual document processing implementation.

- Created new Runner type in document/handlers/runner.go
- Moved test execution logic from DocumentsHandler into Runner
- Updated DocumentsHandler to use Runner for test execution
- Improved error handling and logging

# Enhanced Document Processing Runner

Added actual prompt rendering and simulated LLM response in the RunTest method to better represent the real workflow.

- Updated RunTest to fetch and use the actual prompt template
- Added document content retrieval for processing
- Implemented debug logging for prompt details
- Added simulated LLM response with YAML formatting
- Enhanced test results with more realistic execution metrics

# Enhanced Document Processing Runner YAML Format

Updated the simulated LLM response to use proper YAML code block formatting.

- Modified simulated response to use ```yaml code block format
- Improved readability of extracted document data

# Enhanced YAML Extraction from LLM Response

Added proper markdown code block parsing for LLM responses using the markdown helper package.

- Integrated markdown.ExtractAllBlocks for parsing LLM response
- Added structured YAML content extraction per document
- Improved document content handling with proper indentation

# Enhanced YAML Extraction Logging

Added comprehensive debug logging throughout the YAML extraction process.

- Added logging for markdown block extraction
- Added detailed logging for YAML content processing
- Added document section boundary logging
- Added content extraction result logging

# Improved YAML Block Processing

Changed YAML extraction to process each markdown block independently.

- Modified extraction to iterate over markdown blocks instead of input documents
- Added document map for efficient ID lookups
- Improved empty line handling for document separation
- Added warning when no documents are extracted

# Enhanced YAML Validation and Document Matching

Added proper YAML parsing and validation with improved document matching.

- Added YAML parsing validation using yaml.v3
- Added handling for invalid YAML content
- Added document matching based on ID and name
- Added tracking of unmatched documents
- Improved error reporting in extracted documents

# Enhanced Multi-document YAML Processing

Added support for parsing YAML blocks containing multiple documents.

- Added parseMultiDocYAML helper function
- Updated block processing to handle multiple YAML documents
- Added proper document separation and formatting
- Improved logging for YAML document discovery

# Refactored Document Extraction Logic

Extracted YAML document parsing into a dedicated function for better code organization.

- Created parseExtractedDocuments function for response parsing
- Simplified RunTest method by removing parsing logic
- Added proper error handling for document parsing
- Maintained comprehensive logging

# Use Step Factory for LLM Initialization

Switched from hardcoded OpenAI step to using the standard step factory to support multiple LLM providers.

- Updated LLMHelper to use StandardStepFactory for step creation
- Allows configuration of different LLM providers through settings

# Store Test Results in Prompt Executions

Enhance test result persistence by storing them as prompt executions. This provides:
- Immutable snapshots of prompts at test time
- Consistent storage mechanism using existing prompt execution infrastructure
- Better test history tracking with full execution details

- Updated documents-handler to store test results as prompt executions
- Modified test history retrieval to load from stored executions
- Added prompt snapshot storage in execution metadata

# Add TestResults Conversion Methods

Added helper methods to convert between TestResults and PromptExecution for better data consistency and easier storage:

- Added ToExecution method to convert TestResults to PromptExecution
- Added FromExecution method to convert PromptExecution back to TestResults
- Store complete test data including raw LLM output and rendered prompts
- Improved error handling for data conversion

# Improve TestResults Type Conversion

Enhanced FromExecution to properly handle complex type conversion:

- Added YAML-based serialization/deserialization for Document and ExtractedDocument types
- Improved error handling with detailed error messages
- Maintained type safety when converting between interface{} and concrete types

# Add Document Content to Test Results

Enhanced test results to include original document content:

- Added DocumentContents field to TestResults to store input document content
- Updated test results template to display original document content
- Modified prompt execution storage to include document content
- Improved test result display with input/output comparison

# Limit Document Content Display

Added content truncation for better UI performance and readability:

- Limited document content display to first 500 bytes
- Added truncation indicator when content is shortened
- Maintained full content in storage for processing

# Enhanced Document Template Access

Added DocumentWithContent struct to provide document content directly in templates for better prompt generation.

- Added DocumentWithContent struct that embeds Document and includes Content field
- Modified RunTest to use DocumentWithContent for template rendering
- Simplified document content access in templates

# Document Type Consolidation

Simplified document handling by moving DocumentWithContent to document-types.go and consolidating test results structure.

- Moved DocumentWithContent struct to document-types.go
- Removed DocumentContents from TestResults in favor of Documents field with content
- Updated TestResults to use DocumentWithContent consistently
- Simplified document content handling in Runner

# Add Test Results Identification

Added ID and PromptID fields to TestResults for better tracking and correlation.

- Added ID field to TestResults for unique identification
- Added PromptID field to TestResults for prompt correlation
- Updated ToExecution to use TestResults ID instead of generating new one
- Updated FromExecution to properly copy ID and PromptID fields

# Enhanced Test Results Display

Added test result identification information to the UI for better tracking.

- Added ID and PromptID display in test results header
- Added ID and PromptID to execution details section
- Updated document content display to use DocumentWithContent structure
- Improved content truncation display

# Improve Test Results UI Organization

Enhanced the test results display with collapsible sections for better readability.

- Added collapsible sections for system prompt, template, and LLM response
- Added expand/collapse indicators
- Improved visual hierarchy of test results display

# Add JSON Support to Document Extraction

Added support for parsing JSON code blocks in LLM responses, allowing the system to handle both YAML and JSON formatted documents. This increases compatibility with different LLM output formats and provides more flexibility in document processing.

- Added JSON parsing support to parseExtractedDocuments
- Support both single JSON objects and arrays of JSON objects
- Maintain consistent error handling between YAML and JSON formats
- Use pretty printing for JSON output

# Multiple Generated Queries Support

Added support for storing and displaying multiple generated queries in test results instead of just one query. This allows for better representation of LLM responses that generate multiple valid query variations.

- Updated TestResults struct to store multiple generated queries
- Modified ToExecution and FromExecution methods to handle multiple queries
- Updated Runner to store all parsed queries instead of just the first one
- Maintained backward compatibility with existing data

# Synthesis Testing Implementation with Elasticsearch Documents

Updated synthesis testing to work with Elasticsearch search results instead of file content.

- Modified TestResults to use MatchedDocument type with Elasticsearch content
- Updated test results template to display Elasticsearch document content
- Maintained ToExecution and FromExecution methods for proper data conversion

# Synthesis Handler Mock Implementation

Added mock implementation of synthesis testing with Elasticsearch-like documents.

- Updated Runner to accept MatchedDocument list instead of document IDs
- Added mock document generation in synthesis handler
- Removed document manager dependency from synthesis testing
- Simplified template data to use matched documents directly

# JavaScript Embeddings Wrapper

Added a Goja JavaScript wrapper for the embeddings functionality to enable embedding generation from JavaScript code.

- Created JSEmbeddingsWrapper struct to manage JS runtime and provider
- Implemented generateEmbedding method to create embeddings from text
- Added getModel method to retrieve model information
- Provided RegisterEmbeddings function for easy setup

# JavaScript Step Integration

Added JavaScript bindings for the Step abstraction to enable step execution from JavaScript code.

- Created Promise-based API for async execution
- Added blocking synchronous API
- Implemented callback-based streaming API with cancellation
- Added helper for creating JavaScript LambdaSteps
- Provided example with embeddings step integration

# JavaScript Steps Integration Test

Added a test program to demonstrate and verify the JavaScript Steps integration. The program shows how to:
- Set up a JavaScript runtime with Steps
- Register Go Steps for JavaScript use
- Use all three Step APIs (Promise, Blocking, Callbacks)
- Handle logging and error cases

# JavaScript Promise Integration Documentation

Added documentation for using Steps with JavaScript Promises through goja_nodejs eventloop:

- Added Promise-based execution examples
- Added event loop usage guidelines
- Added cancellation handling examples
- Added error handling best practices

# JavaScript Event Loop Integration

Updated JavaScript test program to use goja_nodejs eventloop for proper Promise handling:

- Added eventloop initialization and cleanup
- Updated test program to run async operations properly
- Added sequential test execution with async/await
- Improved error handling in JavaScript code
- Added proper console.error implementation

# JavaScript Promise Resolution Fix

Updated JavaScript Step wrapper to properly handle Promise resolution using the event loop:

- Added event loop to JSStepWrapper struct
- Modified Promise resolution to use RunOnLoop for all callbacks
- Updated callback handling to use event loop for all JS calls
- Fixed potential race conditions in Promise resolution
- Updated RegisterStep signature to require event loop parameter

# JavaScript Step Documentation Update

Updated steps-js.md to match actual implementation:

- Added event loop integration details
- Updated API examples with actual usage patterns
- Added implementation details section
- Improved best practices for event loop usage
- Removed duplicate Promise execution content

# JavaScript Conversation Wrapper

Added a JavaScript wrapper for the conversation package to enable managing conversations from JavaScript code. The wrapper provides a clean API for creating conversations, adding messages, tool uses, and tool results, and converting back to Go conversation objects.

- Added `conversation-js.go` implementing the JS wrapper
- Added `conversation-js.md` with usage documentation and API reference
- Supports chat messages, tool uses, and tool results
- Provides methods for getting messages and single prompt format
- Enables conversion back to Go conversation objects

# Improved Message Serialization in JS Wrapper

Why: To provide a cleaner, more JavaScript-friendly representation of conversation messages

- Modified GetMessages to return a proper JavaScript array of message objects
- Added type-specific fields for chat messages, tool use, and tool results
- Properly handled image arrays in chat messages
- Unmarshaled tool use input JSON into JavaScript objects

# Fixed JavaScript Object Creation in Conversation Wrapper

Why: To properly create JavaScript objects using Goja instead of relying on automatic map conversion

- Modified GetMessages to use Goja's NewObject and NewArray
- Properly set object properties using Set method
- Fixed array indexing for images and messages
- Ensured proper JavaScript object creation for nested structures

# Extended Conversation JS API

Why: To expose more functionality from the conversation package to JavaScript

- Added message options support to addMessage (metadata, parentID, time, id)
- Added addMessageWithImage for creating messages with attached images
- Added getMessageView for formatted message representations
- Added updateMetadata for modifying message metadata
- Updated documentation with new API methods and examples

# Add combined semantic and keyword search for insights

Added a new search-insights command that combines embedding-based semantic search with keyword/fuzzy matching for better search results.

- Created search-insights.yaml query combining KNN and fuzzy matching
- Added boost factor of 4 to keyword matches
- Includes AUTO fuzziness for better matching
