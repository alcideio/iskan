/*
 * InsightVM Cloud API
 *
 * # Overview   This guide documents the InsightVM Cloud Application Programming Interface (API). This API supports the Representation State Transfer (REST) design pattern. See [Insight Platform API Overview](https://insight.help.rapid7.com/docs/api-overview) for an overview of all Insight Platform APIs.  Versioning is specified in the URL and the base path of this API is:   `https://{region}.api.insight.rapid7.com/vm/{version}/`  Version numbers are numerical and prefixed with the letter `\"v\"`, such as `\"v1\"`.  The region indicates the geo-location of the Insight Platform desired:   | Code  | Region                |  |-------|-----------------------|  | us    | United States         |  | eu    | Europe                |  | ca    | Canada                |  | au    | Australia             |  | ap    | Japan                 |  ## Authorization  Authorization requires a token header `X-Api-Key` and can be generated from the [Insight Platform](https://insight.rapid7.com) key management page. See [Insight Platform API Key](https://insight.help.rapid7.com/docs/managing-platform-api-keys) for more details.  ## Media Types  Unless noted otherwise this API accepts and produces the `application/json` and `application/xml` media types.  Unless otherwise indicated, the default request body media type is `application/json`.   ## Discoverability  All resources respond to the `OPTIONS` request, which allows discoverability of available operations that are supported.  The `OPTIONS` response returns the acceptable HTTP operations on that resource within the `Allow` header. The response is always a `200 OK` status.  ## Verbs and Responses  The following HTTP operations are supported throughout this API. The general usage of the operation and both its failure and success status codes are outlined below.    | <div style=\"width: 70px\">Verb</div>      | Usage                                                                                 | Success     | Failure                                                        | | --------- | ------------------------------------------------------------------------------------- | ----------- | -------------------------------------------------------------- | | `GET`     | Used to retrieve a resource by identifier, or a collection of resources by type.      | `200`       | `400`, `401`, `402`, `404`, `405`, `408`, `410`, `415`, `500`  | | `POST`    | Creates a resource with an application-specified identifier.                          | `201`       | `400`, `401`, `404`, `405`, `408`, `413`, `415`, `500`         | | `POST`    | Performs a request to queue an asynchronous job.                                      | `202`       | `400`, `401`, `405`, `408`, `410`, `413`, `415`, `500`         | | `PUT`     | Creates a resource with a client-specified identifier.                                | `200`       | `400`, `401`, `403`, `405`, `408`, `410`, `413`, `415`, `500`  | | `PUT`     | Performs a full update of a resource with a specified identifier.                     | `201`       | `400`, `401`, `403`, `405`, `408`, `410`, `413`, `415`, `500`  | | `DELETE`  | Deletes a resource by identifier or an entire collection of resources.                | `204`       | `400`, `401`, `405`, `408`, `410`, `413`, `415`, `500`         | | `OPTIONS` | Requests what operations are available on a resource.                                 | `200`       | `401`, `404`, `405`, `408`, `500`                              |  ## Resources  Resource names represent nouns and identify the entity being manipulated or accessed. All collection resources are  pluralized to indicate to the client they are interacting with a collection of multiple resources of the same type. Singular resource names are used when there exists only one resource available to interact with.  The following naming conventions are used by this API:  | Type                                          | Case                     | | --------------------------------------------- | ------------------------ | | Resource names                                | `strike-through-case`    | | Header, body, and query parameters parameters | `camelCase`              | | JSON fields and property names                | `snake_case`             |  ### Collections  A collection resource is a parent resource for instance resources, but can itself be retrieved and operated on  independently. Collection resources use a pluralized resource name. The resource path for collection resources follow  the convention:  ``` /{resource_name} ```  Collection resources can support the `GET`, `POST`, `PUT`, and `DELETE` operations.  #### GET  The `GET` operation invoked on a collection resource indicates a request to retrieve all, or some, of the entities  contained within the collection. This also includes the optional capability to filter or search resources during the request. The response from a collection listing is a paginated document.  #### POST  The `POST` is a non-idempotent operation that allows for the creation of a new resource when the resource identifier  is not provided by the system during the creation operation (i.e. the Security Console generates the identifier). The content of the `POST` request is sent in the request body. The response to a successful `POST` request should be a  `201 CREATED` with a valid `Location` header field set to the URI that can be used to access to the newly  created resource.   The `POST` to a collection resource can also be used to interact with asynchronous resources. In this situation,  instead of a `201 CREATED` response, the `202 ACCEPTED` response indicates that processing of the request is not fully  complete but has been accepted for future processing. This request will respond similarly with a `Location` header with  link to the job-oriented asynchronous resource that was created and/or queued.  #### PUT  The `PUT` is an idempotent operation that either performs a create with user-supplied identity, or a full replace or update of a resource by a known identifier. The response to a `PUT` operation to create an entity is a `201 Created` with a valid `Location` header field set to the URI that can be used to access to the newly created resource.  `PUT` on a collection resource replaces all values in the collection. The typical response to a `PUT` operation that  updates an entity is hypermedia links, which may link to related resources caused by the side-effects of the changes  performed.  #### DELETE  The `DELETE` is an idempotent operation that physically deletes a resource, or removes an association between resources. The typical response to a `DELETE` operation is hypermedia links, which may link to related resources caused by the  side-effects of the changes performed.  ### Instances  An instance resource is a \"leaf\" level resource that may be retrieved, optionally nested within a collection resource. Instance resources are usually retrievable with opaque identifiers. The resource path for instance resources follows  the convention:  ``` /{resource_name}/{instance_id}... ```  Instance resources can support the `GET`, `PUT`, `POST`, `PATCH` and `DELETE` operations.  #### GET  Retrieves the details of a specific resource by its identifier. The details retrieved can be controlled through  property selection and property views. The content of the resource is returned within the body of the response in the  acceptable media type.   #### PUT  Allows for and idempotent \"full update\" (complete replacement) on a specific resource. If the resource does not exist,  it will be created; if it does exist, it is completely overwritten. Any omitted properties in the request are assumed to  be undefined/null. For \"partial updates\" use `POST` or `PATCH` instead.   The content of the `PUT` request is sent in the request body. The identifier of the resource is specified within the URL  (not the request body). The response to a successful `PUT` request is a `201 CREATED` to represent the created status,  with a valid `Location` header field set to the URI that can be used to access to the newly created (or fully replaced)  resource.   #### POST  Performs a non-idempotent creation of a new resource. The `POST` of an instance resource most commonly occurs with the  use of nested resources (e.g. searching on a parent collection resource). The response to a `POST` of an instance  resource is typically a `200 OK` if the resource is non-persistent, and a `201 CREATED` if there is a resource  created/persisted as a result of the operation. This varies by endpoint.  #### PATCH  The `PATCH` operation is used to perform a partial update of a resource. `PATCH` is a non-idempotent operation that enforces an atomic mutation of a resource. Only the properties specified in the request are to be overwritten on the  resource it is applied to. If a property is missing, it is assumed to not have changed.  #### DELETE  Permanently removes the individual resource from the system. If the resource is an association between resources, only  the association is removed, not the resources themselves. A successful deletion of the resource should return  `204 NO CONTENT` with no response body. This operation is not fully idempotent, as follow-up requests to delete a  non-existent resource should return a `404 NOT FOUND`.  ## Formats  ### Dates & Times  Dates and/or times are specified as strings in the ISO 8601 format(s). The following formats are supported as input:  | Value                       | Format                                                 | Notes                                                 | | --------------------------- | ------------------------------------------------------ | ----------------------------------------------------- | | Date                        | YYYY-MM-DD                                             | Defaults to 12 am UTC (if used for a date & time      | | Date & time only            | YYYY-MM-DD'T'hh:mm:ss[.nnn]                            | Defaults to UTC                                       | | Date & time in UTC          | YYYY-MM-DD'T'hh:mm:ss[.nnn]Z                           |                                                       | | Date & time w/ offset       | YYYY-MM-DD'T'hh:mm:ss[.nnn][+&#124;-]hh:mm             |                                                       | | Date & time w/ zone-offset  | YYYY-MM-DD'T'hh:mm:ss[.nnn][+&#124;-]hh:mm[<zone-id>]  |                                                       |   ### Timezones  Timezones are specified in the regional zone format, such as `\"America/Los_Angeles\"`, `\"Asia/Tokyo\"`, or `\"GMT\"`.   ### Paging  Pagination may be supported on collection resources using a combination of two query parameters, `page` and `size`.  The page parameter dictates the  zero-based index of the page to retrieve, and the `size` indicates the size of the page.   For example, `/resources?page=2&size=10` will return page 3, with 10 records per page, giving results 21-30.  The maximum page size for a request is 1000.  Some paginated endpoints may supported \"cursored\" pages, allowing for a guaranteed consistent view of data across page boundaries. Cursored page requests support a consistent, sequential way to access data across pages. Only if this option  is used are you guaranteed that you will read a record once and only once in any page (\"repeatable read\").  If not supported, or not specified, the results may shift across page boundaries while they are being read as data updates  (\"read committed\"). The `cursor` property is used to follow the same chain of paginated requests from page to page. Each  request will change the value of the next cursor to use on the subsequent page, and may only be used to iterate sequentially through pages.  The response to a paginated request follows the format:  ```json {    data\": [        ...     ],    \"metadata\": {        \"index\": ...,       \"size\": ...,       \"sort\": ...,       \"total_data\": ...,       \"total_pages\": ...,       \"cursor\": ...    },    \"links\": [        \"first\" : {          \"href\" : \"...\"        },        \"prev\" : {          \"href\" : \"...\"        },        \"self\" : {          \"href\" : \"...\"        },        \"next\" : {          \"href\" : \"...\"        },        \"last\" : {          \"href\" : \"...\"        }     ] } ```  The `data` property is an array of the resources being retrieved from the endpoint, each which should contain at  minimum a \"self\" relation hypermedia link. The `metadata` property outlines the details of the current page and total possible pages. The object for the page includes the following properties:  - `index` - The page number (zero-based) of the page returned. - `size` - The size of the pages, which is less than or equal to the maximum page size. - `total_data` - The total amount of resources available across all pages. - `total_pages` - The total amount of pages. - `cursor` - An optional cursor for \"cursored\" page requests  The last property of the paged response is the `links` array, which contains all available hypermedia links. For  paginated responses, the \"self\", \"next\", \"previous\", \"first\", and \"last\" links are returned. The \"self\" link must always be returned and should contain a link to allow the client to replicate the original request against the  collection resource in an identical manner to that in which it was invoked.   The \"next\" and \"previous\" links are present if either or both there exists a previous or next page, respectively.  The \"next\" and \"previous\" links have hrefs that allow \"natural movement\" to the next page, that is all parameters  required to move the next page are provided in the link. The \"first\" and \"last\" links provide references to the first and last pages respectively. If the page is \"cursored\" the cursor is automatically incorporated into the pagination links.  ### Sorting  Sorting is supported on paginated resources with the `sort` query parameter(s). The sort query parameter(s) supports  identifying a single or multi-property sort with a single or multi-direction output. The format of the parameter is:  ``` sort=property[,ASC|DESC]... ```  Therefore, the request `/resources?sort=name,title,DESC` would return the results sorted by the name and title  descending, in that order. The sort directions are either ascending `ASC` or descending `DESC`. With single-order  sorting, all properties are sorted in the same direction. To sort the results with varying orders by property,  multiple sort parameters are passed.    For example, the request `/resources?sort=name,ASC&sort=title,DESC` would sort by name ascending and title  descending, in that order.  ## Responses  The following response statuses may be returned by this API.     | Status | Meaning                  | Usage                                                                                                                                                                    | | ------ | ------------------------ |------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | | `200`  | OK                       | The operation performed without error according to the specification of the request, and no more specific 2xx code is suitable.                                          | | `201`  | Created                  | A create request has been fulfilled and a resource has been created. The resource is available as the URI specified in the response, including the `Location` header.    | | `202`  | Accepted                 | An asynchronous task has been accepted, but not guaranteed, to be processed in the future.                                                                               | | `400`  | Bad Request              | The request was invalid or cannot be otherwise served. The request is not likely to succeed in the future without modifications.                                         | | `401`  | Unauthorized             | The user is unauthorized to perform the operation requested, or does not maintain permissions to perform the operation on the resource specified.                        | | `403`  | Forbidden                | The resource exists to which the user has access, but the operating requested is not permitted.                                                                          | | `404`  | Not Found                | The resource specified could not be located, does not exist, or an unauthenticated client does not have permissions to a resource.                                       | | `405`  | Method Not Allowed       | The operations may not be performed on the specific resource. Allowed operations are returned and may be performed on the resource.                                      | | `408`  | Request Timeout          | The client has failed to complete a request in a timely manner and the request has been discarded.                                                                       | | `413`  | Request Entity Too Large | The request being provided is too large for the server to accept processing.                                                                                             | | `415`  | Unsupported Media Type   | The media type is not supported for the requested resource.                                                                                                              | | `500`  | Internal Server Error    | An internal and unexpected error has occurred on the server at no fault of the client.                                                                                   |  ### Errors  Any error responses can provide a response body with a message to the client indicating more information (if applicable)  to aid debugging of the error. All 4xx and 5xx responses will return an error response in the body. The format of the  response is as follows:  ```json {    \"status\": <statusCode>,    \"message\": <message>,    \"localized_message\": <message>,    \"links\" : [ {       \"rel\" : \"...\",       \"href\" : \"...\"     } ] }   ```   The `status` property is the same as the HTTP status returned in the response, to ease client parsing. The message  property is a localized message in the request client's locale (if applicable) that articulates the nature of the  error. The last property is the `links` property.  ### Security  The response statuses 401, 403 and 404 need special consideration for security purposes. As necessary,  error statuses and messages may be obscured to strengthen security and prevent information exposure. The following is a  guideline for privileged resource response statuses:  | Use Case                                                           | Access             | Resource           | Permission   | Status       | | ------------------------------------------------------------------ | ------------------ |------------------- | ------------ | ------------ | | Unauthenticated access to an unauthenticated resource.             | Unauthenticated    | Unauthenticated    | Yes          | `20x`        | | Unauthenticated access to an authenticated resource.               | Unauthenticated    | Authenticated      | No           | `401`        | | Unauthenticated access to an authenticated resource.               | Unauthenticated    | Non-existent       | No           | `401`        | | Authenticated access to a unauthenticated resource.                | Authenticated      | Unauthenticated    | Yes          | `20x`        | | Authenticated access to an authenticated, unprivileged resource.   | Authenticated      | Authenticated      | No           | `404`        | | Authenticated access to an authenticated, privileged resource.     | Authenticated      | Authenticated      | Yes          | `20x`        | | Authenticated access to an authenticated, non-existent resource    | Authenticated      | Non-existent       | Yes          | `404`        |  ### Headers  Commonly used response headers include:  | Header                     |  Example                          | Purpose                                                         | | -------------------------- | --------------------------------- | --------------------------------------------------------------- | | `Allow`                    | `OPTIONS, GET`                    | Defines the allowable HTTP operations on a resource.            | | `Cache-Control`            | `no-store, must-revalidate`       | Disables caching of resources (as they are all dynamic).        | | `Content-Encoding`         | `gzip`                            | The encoding of the response body (if any).                     | | `Location`                 |                                   | Refers to the URI of the resource created by a request.         | | `Transfer-Encoding`        | `chunked`                         | Specified the encoding used to transform response.              | | `Retry-After`              | 5000                              | Indicates the time to wait before retrying a request.           | | `X-Content-Type-Options`   | `nosniff`                         | Disables MIME type sniffing.                                    | | `X-XSS-Protection`         | `1; mode=block`                   | Enables XSS filter protection.                                  | | `X-Frame-Options`          | `SAMEORIGIN`                      | Prevents rendering in a frame from a different origin.          | | `X-UA-Compatible`          | `IE=edge,chrome=1`                | Specifies the browser mode to render in.                        |  ### Format  When `application/json` is returned in the response body it is always pretty-printed (indented, human readable output).  Additionally, gzip compression/encoding is supported on all responses.   #### Dates & Times  Dates or times are returned as strings in the ISO 8601 'extended' format. When a date and time is returned (instant) the value is converted to UTC.  For example:  | Value           | Format                         | Example               | | --------------- | ------------------------------ | --------------------- | | Date            | `YYYY-MM-DD`                   | 2017-12-03            | | Date & Time     | `YYYY-MM-DD'T'hh:mm:ss[.nnn]Z` | 2017-12-03T10:15:30Z  |  # Authentication  <!-- ReDoc-Inject: <security-definitions> -->
 *
 * API version: 1.0.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// Layer ${layer.description}
type Layer struct {
	Assessment *ImageAssessment `json:"assessment,omitempty"`
	// The name of the author of the layer.
	Author  *string            `json:"author,omitempty"`
	Command *DockerfileCommand `json:"command,omitempty"`
	// Commands that were executed to create the layer.
	Commands *string `json:"commands,omitempty"`
	// Textual comment attached to the layer.
	Comment *string `json:"comment,omitempty"`
	// The ISO instant created time.
	Created *string `json:"created,omitempty"`
	// Whether the layer is empty or not.
	Empty *bool `json:"empty,omitempty"`
	// The identifier of the layer, in hash format `[<algorithm>:]<hash>`.
	Id              *string               `json:"id,omitempty"`
	OperatingSystem *ImageOperatingSystem `json:"operating_system,omitempty"`
	// Number of packages in the layer.
	PackageCount *int32 `json:"package_count,omitempty"`
	// Installed packages detected in the layer.
	Packages *[]Package `json:"packages,omitempty"`
	// The identifier of the parent layer on top of which the current layer is built, in hash format `[<algorithm>:]<hash>`.
	ParentId *string `json:"parent_id,omitempty"`
	// Position of the layer (also called intermediate images) during the Docker image build.
	Position *int32 `json:"position,omitempty"`
	// Size of the image layer in bytes.
	Size *int64 `json:"size,omitempty"`
}

// NewLayer instantiates a new Layer object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewLayer() *Layer {
	this := Layer{}
	return &this
}

// NewLayerWithDefaults instantiates a new Layer object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewLayerWithDefaults() *Layer {
	this := Layer{}
	return &this
}

// GetAssessment returns the Assessment field value if set, zero value otherwise.
func (o *Layer) GetAssessment() ImageAssessment {
	if o == nil || o.Assessment == nil {
		var ret ImageAssessment
		return ret
	}
	return *o.Assessment
}

// GetAssessmentOk returns a tuple with the Assessment field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Layer) GetAssessmentOk() (*ImageAssessment, bool) {
	if o == nil || o.Assessment == nil {
		return nil, false
	}
	return o.Assessment, true
}

// HasAssessment returns a boolean if a field has been set.
func (o *Layer) HasAssessment() bool {
	if o != nil && o.Assessment != nil {
		return true
	}

	return false
}

// SetAssessment gets a reference to the given ImageAssessment and assigns it to the Assessment field.
func (o *Layer) SetAssessment(v ImageAssessment) {
	o.Assessment = &v
}

// GetAuthor returns the Author field value if set, zero value otherwise.
func (o *Layer) GetAuthor() string {
	if o == nil || o.Author == nil {
		var ret string
		return ret
	}
	return *o.Author
}

// GetAuthorOk returns a tuple with the Author field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Layer) GetAuthorOk() (*string, bool) {
	if o == nil || o.Author == nil {
		return nil, false
	}
	return o.Author, true
}

// HasAuthor returns a boolean if a field has been set.
func (o *Layer) HasAuthor() bool {
	if o != nil && o.Author != nil {
		return true
	}

	return false
}

// SetAuthor gets a reference to the given string and assigns it to the Author field.
func (o *Layer) SetAuthor(v string) {
	o.Author = &v
}

// GetCommand returns the Command field value if set, zero value otherwise.
func (o *Layer) GetCommand() DockerfileCommand {
	if o == nil || o.Command == nil {
		var ret DockerfileCommand
		return ret
	}
	return *o.Command
}

// GetCommandOk returns a tuple with the Command field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Layer) GetCommandOk() (*DockerfileCommand, bool) {
	if o == nil || o.Command == nil {
		return nil, false
	}
	return o.Command, true
}

// HasCommand returns a boolean if a field has been set.
func (o *Layer) HasCommand() bool {
	if o != nil && o.Command != nil {
		return true
	}

	return false
}

// SetCommand gets a reference to the given DockerfileCommand and assigns it to the Command field.
func (o *Layer) SetCommand(v DockerfileCommand) {
	o.Command = &v
}

// GetCommands returns the Commands field value if set, zero value otherwise.
func (o *Layer) GetCommands() string {
	if o == nil || o.Commands == nil {
		var ret string
		return ret
	}
	return *o.Commands
}

// GetCommandsOk returns a tuple with the Commands field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Layer) GetCommandsOk() (*string, bool) {
	if o == nil || o.Commands == nil {
		return nil, false
	}
	return o.Commands, true
}

// HasCommands returns a boolean if a field has been set.
func (o *Layer) HasCommands() bool {
	if o != nil && o.Commands != nil {
		return true
	}

	return false
}

// SetCommands gets a reference to the given string and assigns it to the Commands field.
func (o *Layer) SetCommands(v string) {
	o.Commands = &v
}

// GetComment returns the Comment field value if set, zero value otherwise.
func (o *Layer) GetComment() string {
	if o == nil || o.Comment == nil {
		var ret string
		return ret
	}
	return *o.Comment
}

// GetCommentOk returns a tuple with the Comment field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Layer) GetCommentOk() (*string, bool) {
	if o == nil || o.Comment == nil {
		return nil, false
	}
	return o.Comment, true
}

// HasComment returns a boolean if a field has been set.
func (o *Layer) HasComment() bool {
	if o != nil && o.Comment != nil {
		return true
	}

	return false
}

// SetComment gets a reference to the given string and assigns it to the Comment field.
func (o *Layer) SetComment(v string) {
	o.Comment = &v
}

// GetCreated returns the Created field value if set, zero value otherwise.
func (o *Layer) GetCreated() string {
	if o == nil || o.Created == nil {
		var ret string
		return ret
	}
	return *o.Created
}

// GetCreatedOk returns a tuple with the Created field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Layer) GetCreatedOk() (*string, bool) {
	if o == nil || o.Created == nil {
		return nil, false
	}
	return o.Created, true
}

// HasCreated returns a boolean if a field has been set.
func (o *Layer) HasCreated() bool {
	if o != nil && o.Created != nil {
		return true
	}

	return false
}

// SetCreated gets a reference to the given string and assigns it to the Created field.
func (o *Layer) SetCreated(v string) {
	o.Created = &v
}

// GetEmpty returns the Empty field value if set, zero value otherwise.
func (o *Layer) GetEmpty() bool {
	if o == nil || o.Empty == nil {
		var ret bool
		return ret
	}
	return *o.Empty
}

// GetEmptyOk returns a tuple with the Empty field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Layer) GetEmptyOk() (*bool, bool) {
	if o == nil || o.Empty == nil {
		return nil, false
	}
	return o.Empty, true
}

// HasEmpty returns a boolean if a field has been set.
func (o *Layer) HasEmpty() bool {
	if o != nil && o.Empty != nil {
		return true
	}

	return false
}

// SetEmpty gets a reference to the given bool and assigns it to the Empty field.
func (o *Layer) SetEmpty(v bool) {
	o.Empty = &v
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *Layer) GetId() string {
	if o == nil || o.Id == nil {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Layer) GetIdOk() (*string, bool) {
	if o == nil || o.Id == nil {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *Layer) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *Layer) SetId(v string) {
	o.Id = &v
}

// GetOperatingSystem returns the OperatingSystem field value if set, zero value otherwise.
func (o *Layer) GetOperatingSystem() ImageOperatingSystem {
	if o == nil || o.OperatingSystem == nil {
		var ret ImageOperatingSystem
		return ret
	}
	return *o.OperatingSystem
}

// GetOperatingSystemOk returns a tuple with the OperatingSystem field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Layer) GetOperatingSystemOk() (*ImageOperatingSystem, bool) {
	if o == nil || o.OperatingSystem == nil {
		return nil, false
	}
	return o.OperatingSystem, true
}

// HasOperatingSystem returns a boolean if a field has been set.
func (o *Layer) HasOperatingSystem() bool {
	if o != nil && o.OperatingSystem != nil {
		return true
	}

	return false
}

// SetOperatingSystem gets a reference to the given ImageOperatingSystem and assigns it to the OperatingSystem field.
func (o *Layer) SetOperatingSystem(v ImageOperatingSystem) {
	o.OperatingSystem = &v
}

// GetPackageCount returns the PackageCount field value if set, zero value otherwise.
func (o *Layer) GetPackageCount() int32 {
	if o == nil || o.PackageCount == nil {
		var ret int32
		return ret
	}
	return *o.PackageCount
}

// GetPackageCountOk returns a tuple with the PackageCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Layer) GetPackageCountOk() (*int32, bool) {
	if o == nil || o.PackageCount == nil {
		return nil, false
	}
	return o.PackageCount, true
}

// HasPackageCount returns a boolean if a field has been set.
func (o *Layer) HasPackageCount() bool {
	if o != nil && o.PackageCount != nil {
		return true
	}

	return false
}

// SetPackageCount gets a reference to the given int32 and assigns it to the PackageCount field.
func (o *Layer) SetPackageCount(v int32) {
	o.PackageCount = &v
}

// GetPackages returns the Packages field value if set, zero value otherwise.
func (o *Layer) GetPackages() []Package {
	if o == nil || o.Packages == nil {
		var ret []Package
		return ret
	}
	return *o.Packages
}

// GetPackagesOk returns a tuple with the Packages field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Layer) GetPackagesOk() (*[]Package, bool) {
	if o == nil || o.Packages == nil {
		return nil, false
	}
	return o.Packages, true
}

// HasPackages returns a boolean if a field has been set.
func (o *Layer) HasPackages() bool {
	if o != nil && o.Packages != nil {
		return true
	}

	return false
}

// SetPackages gets a reference to the given []Package and assigns it to the Packages field.
func (o *Layer) SetPackages(v []Package) {
	o.Packages = &v
}

// GetParentId returns the ParentId field value if set, zero value otherwise.
func (o *Layer) GetParentId() string {
	if o == nil || o.ParentId == nil {
		var ret string
		return ret
	}
	return *o.ParentId
}

// GetParentIdOk returns a tuple with the ParentId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Layer) GetParentIdOk() (*string, bool) {
	if o == nil || o.ParentId == nil {
		return nil, false
	}
	return o.ParentId, true
}

// HasParentId returns a boolean if a field has been set.
func (o *Layer) HasParentId() bool {
	if o != nil && o.ParentId != nil {
		return true
	}

	return false
}

// SetParentId gets a reference to the given string and assigns it to the ParentId field.
func (o *Layer) SetParentId(v string) {
	o.ParentId = &v
}

// GetPosition returns the Position field value if set, zero value otherwise.
func (o *Layer) GetPosition() int32 {
	if o == nil || o.Position == nil {
		var ret int32
		return ret
	}
	return *o.Position
}

// GetPositionOk returns a tuple with the Position field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Layer) GetPositionOk() (*int32, bool) {
	if o == nil || o.Position == nil {
		return nil, false
	}
	return o.Position, true
}

// HasPosition returns a boolean if a field has been set.
func (o *Layer) HasPosition() bool {
	if o != nil && o.Position != nil {
		return true
	}

	return false
}

// SetPosition gets a reference to the given int32 and assigns it to the Position field.
func (o *Layer) SetPosition(v int32) {
	o.Position = &v
}

// GetSize returns the Size field value if set, zero value otherwise.
func (o *Layer) GetSize() int64 {
	if o == nil || o.Size == nil {
		var ret int64
		return ret
	}
	return *o.Size
}

// GetSizeOk returns a tuple with the Size field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Layer) GetSizeOk() (*int64, bool) {
	if o == nil || o.Size == nil {
		return nil, false
	}
	return o.Size, true
}

// HasSize returns a boolean if a field has been set.
func (o *Layer) HasSize() bool {
	if o != nil && o.Size != nil {
		return true
	}

	return false
}

// SetSize gets a reference to the given int64 and assigns it to the Size field.
func (o *Layer) SetSize(v int64) {
	o.Size = &v
}

func (o Layer) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Assessment != nil {
		toSerialize["assessment"] = o.Assessment
	}
	if o.Author != nil {
		toSerialize["author"] = o.Author
	}
	if o.Command != nil {
		toSerialize["command"] = o.Command
	}
	if o.Commands != nil {
		toSerialize["commands"] = o.Commands
	}
	if o.Comment != nil {
		toSerialize["comment"] = o.Comment
	}
	if o.Created != nil {
		toSerialize["created"] = o.Created
	}
	if o.Empty != nil {
		toSerialize["empty"] = o.Empty
	}
	if o.Id != nil {
		toSerialize["id"] = o.Id
	}
	if o.OperatingSystem != nil {
		toSerialize["operating_system"] = o.OperatingSystem
	}
	if o.PackageCount != nil {
		toSerialize["package_count"] = o.PackageCount
	}
	if o.Packages != nil {
		toSerialize["packages"] = o.Packages
	}
	if o.ParentId != nil {
		toSerialize["parent_id"] = o.ParentId
	}
	if o.Position != nil {
		toSerialize["position"] = o.Position
	}
	if o.Size != nil {
		toSerialize["size"] = o.Size
	}
	return json.Marshal(toSerialize)
}

type NullableLayer struct {
	value *Layer
	isSet bool
}

func (v NullableLayer) Get() *Layer {
	return v.value
}

func (v *NullableLayer) Set(val *Layer) {
	v.value = val
	v.isSet = true
}

func (v NullableLayer) IsSet() bool {
	return v.isSet
}

func (v *NullableLayer) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableLayer(val *Layer) *NullableLayer {
	return &NullableLayer{value: val, isSet: true}
}

func (v NullableLayer) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableLayer) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
