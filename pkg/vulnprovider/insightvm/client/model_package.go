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

// Package ${package.description}
type Package struct {
	Assessment *ImageAssessment `json:"assessment,omitempty"`
	// Provides description about the package.
	Description *string `json:"description,omitempty"`
	// Epoch time when the package was created.
	Epoch *string `json:"epoch,omitempty"`
	// URL link to the home page of the package.
	HomePage *string `json:"home_page,omitempty"`
	// Package id.
	Id *int64 `json:"id,omitempty"`
	// Describes the type of licensing model.
	License *string `json:"license,omitempty"`
	// Contact details of the maintainer of the package.
	Maintainer *string `json:"maintainer,omitempty"`
	// Name of the package.
	Name *string `json:"name,omitempty"`
	// Operating system architecture type.
	OsArchitecture *string `json:"os_architecture,omitempty"`
	// Describes the type of the operating system family. Some examples include `Windows`, `Linux`, `MacOs` etc.
	OsFamily *string `json:"os_family,omitempty"`
	// Name of the operating system.
	OsName *string `json:"os_name,omitempty"`
	// Operating system vendor.
	OsVendor *string `json:"os_vendor,omitempty"`
	// Package operating system version.
	OsVersion *string `json:"os_version,omitempty"`
	// Package release information.
	Release *string `json:"release,omitempty"`
	// Size of the package in bytes.
	Size *int64 `json:"size,omitempty"`
	// Source of the package.
	Source *string `json:"source,omitempty"`
	// Package type.
	Type *string `json:"type,omitempty"`
	// Version of the package.
	Version *string `json:"version,omitempty"`
}

// NewPackage instantiates a new Package object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPackage() *Package {
	this := Package{}
	return &this
}

// NewPackageWithDefaults instantiates a new Package object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPackageWithDefaults() *Package {
	this := Package{}
	return &this
}

// GetAssessment returns the Assessment field value if set, zero value otherwise.
func (o *Package) GetAssessment() ImageAssessment {
	if o == nil || o.Assessment == nil {
		var ret ImageAssessment
		return ret
	}
	return *o.Assessment
}

// GetAssessmentOk returns a tuple with the Assessment field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Package) GetAssessmentOk() (*ImageAssessment, bool) {
	if o == nil || o.Assessment == nil {
		return nil, false
	}
	return o.Assessment, true
}

// HasAssessment returns a boolean if a field has been set.
func (o *Package) HasAssessment() bool {
	if o != nil && o.Assessment != nil {
		return true
	}

	return false
}

// SetAssessment gets a reference to the given ImageAssessment and assigns it to the Assessment field.
func (o *Package) SetAssessment(v ImageAssessment) {
	o.Assessment = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *Package) GetDescription() string {
	if o == nil || o.Description == nil {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Package) GetDescriptionOk() (*string, bool) {
	if o == nil || o.Description == nil {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *Package) HasDescription() bool {
	if o != nil && o.Description != nil {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *Package) SetDescription(v string) {
	o.Description = &v
}

// GetEpoch returns the Epoch field value if set, zero value otherwise.
func (o *Package) GetEpoch() string {
	if o == nil || o.Epoch == nil {
		var ret string
		return ret
	}
	return *o.Epoch
}

// GetEpochOk returns a tuple with the Epoch field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Package) GetEpochOk() (*string, bool) {
	if o == nil || o.Epoch == nil {
		return nil, false
	}
	return o.Epoch, true
}

// HasEpoch returns a boolean if a field has been set.
func (o *Package) HasEpoch() bool {
	if o != nil && o.Epoch != nil {
		return true
	}

	return false
}

// SetEpoch gets a reference to the given string and assigns it to the Epoch field.
func (o *Package) SetEpoch(v string) {
	o.Epoch = &v
}

// GetHomePage returns the HomePage field value if set, zero value otherwise.
func (o *Package) GetHomePage() string {
	if o == nil || o.HomePage == nil {
		var ret string
		return ret
	}
	return *o.HomePage
}

// GetHomePageOk returns a tuple with the HomePage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Package) GetHomePageOk() (*string, bool) {
	if o == nil || o.HomePage == nil {
		return nil, false
	}
	return o.HomePage, true
}

// HasHomePage returns a boolean if a field has been set.
func (o *Package) HasHomePage() bool {
	if o != nil && o.HomePage != nil {
		return true
	}

	return false
}

// SetHomePage gets a reference to the given string and assigns it to the HomePage field.
func (o *Package) SetHomePage(v string) {
	o.HomePage = &v
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *Package) GetId() int64 {
	if o == nil || o.Id == nil {
		var ret int64
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Package) GetIdOk() (*int64, bool) {
	if o == nil || o.Id == nil {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *Package) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// SetId gets a reference to the given int64 and assigns it to the Id field.
func (o *Package) SetId(v int64) {
	o.Id = &v
}

// GetLicense returns the License field value if set, zero value otherwise.
func (o *Package) GetLicense() string {
	if o == nil || o.License == nil {
		var ret string
		return ret
	}
	return *o.License
}

// GetLicenseOk returns a tuple with the License field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Package) GetLicenseOk() (*string, bool) {
	if o == nil || o.License == nil {
		return nil, false
	}
	return o.License, true
}

// HasLicense returns a boolean if a field has been set.
func (o *Package) HasLicense() bool {
	if o != nil && o.License != nil {
		return true
	}

	return false
}

// SetLicense gets a reference to the given string and assigns it to the License field.
func (o *Package) SetLicense(v string) {
	o.License = &v
}

// GetMaintainer returns the Maintainer field value if set, zero value otherwise.
func (o *Package) GetMaintainer() string {
	if o == nil || o.Maintainer == nil {
		var ret string
		return ret
	}
	return *o.Maintainer
}

// GetMaintainerOk returns a tuple with the Maintainer field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Package) GetMaintainerOk() (*string, bool) {
	if o == nil || o.Maintainer == nil {
		return nil, false
	}
	return o.Maintainer, true
}

// HasMaintainer returns a boolean if a field has been set.
func (o *Package) HasMaintainer() bool {
	if o != nil && o.Maintainer != nil {
		return true
	}

	return false
}

// SetMaintainer gets a reference to the given string and assigns it to the Maintainer field.
func (o *Package) SetMaintainer(v string) {
	o.Maintainer = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *Package) GetName() string {
	if o == nil || o.Name == nil {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Package) GetNameOk() (*string, bool) {
	if o == nil || o.Name == nil {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *Package) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *Package) SetName(v string) {
	o.Name = &v
}

// GetOsArchitecture returns the OsArchitecture field value if set, zero value otherwise.
func (o *Package) GetOsArchitecture() string {
	if o == nil || o.OsArchitecture == nil {
		var ret string
		return ret
	}
	return *o.OsArchitecture
}

// GetOsArchitectureOk returns a tuple with the OsArchitecture field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Package) GetOsArchitectureOk() (*string, bool) {
	if o == nil || o.OsArchitecture == nil {
		return nil, false
	}
	return o.OsArchitecture, true
}

// HasOsArchitecture returns a boolean if a field has been set.
func (o *Package) HasOsArchitecture() bool {
	if o != nil && o.OsArchitecture != nil {
		return true
	}

	return false
}

// SetOsArchitecture gets a reference to the given string and assigns it to the OsArchitecture field.
func (o *Package) SetOsArchitecture(v string) {
	o.OsArchitecture = &v
}

// GetOsFamily returns the OsFamily field value if set, zero value otherwise.
func (o *Package) GetOsFamily() string {
	if o == nil || o.OsFamily == nil {
		var ret string
		return ret
	}
	return *o.OsFamily
}

// GetOsFamilyOk returns a tuple with the OsFamily field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Package) GetOsFamilyOk() (*string, bool) {
	if o == nil || o.OsFamily == nil {
		return nil, false
	}
	return o.OsFamily, true
}

// HasOsFamily returns a boolean if a field has been set.
func (o *Package) HasOsFamily() bool {
	if o != nil && o.OsFamily != nil {
		return true
	}

	return false
}

// SetOsFamily gets a reference to the given string and assigns it to the OsFamily field.
func (o *Package) SetOsFamily(v string) {
	o.OsFamily = &v
}

// GetOsName returns the OsName field value if set, zero value otherwise.
func (o *Package) GetOsName() string {
	if o == nil || o.OsName == nil {
		var ret string
		return ret
	}
	return *o.OsName
}

// GetOsNameOk returns a tuple with the OsName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Package) GetOsNameOk() (*string, bool) {
	if o == nil || o.OsName == nil {
		return nil, false
	}
	return o.OsName, true
}

// HasOsName returns a boolean if a field has been set.
func (o *Package) HasOsName() bool {
	if o != nil && o.OsName != nil {
		return true
	}

	return false
}

// SetOsName gets a reference to the given string and assigns it to the OsName field.
func (o *Package) SetOsName(v string) {
	o.OsName = &v
}

// GetOsVendor returns the OsVendor field value if set, zero value otherwise.
func (o *Package) GetOsVendor() string {
	if o == nil || o.OsVendor == nil {
		var ret string
		return ret
	}
	return *o.OsVendor
}

// GetOsVendorOk returns a tuple with the OsVendor field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Package) GetOsVendorOk() (*string, bool) {
	if o == nil || o.OsVendor == nil {
		return nil, false
	}
	return o.OsVendor, true
}

// HasOsVendor returns a boolean if a field has been set.
func (o *Package) HasOsVendor() bool {
	if o != nil && o.OsVendor != nil {
		return true
	}

	return false
}

// SetOsVendor gets a reference to the given string and assigns it to the OsVendor field.
func (o *Package) SetOsVendor(v string) {
	o.OsVendor = &v
}

// GetOsVersion returns the OsVersion field value if set, zero value otherwise.
func (o *Package) GetOsVersion() string {
	if o == nil || o.OsVersion == nil {
		var ret string
		return ret
	}
	return *o.OsVersion
}

// GetOsVersionOk returns a tuple with the OsVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Package) GetOsVersionOk() (*string, bool) {
	if o == nil || o.OsVersion == nil {
		return nil, false
	}
	return o.OsVersion, true
}

// HasOsVersion returns a boolean if a field has been set.
func (o *Package) HasOsVersion() bool {
	if o != nil && o.OsVersion != nil {
		return true
	}

	return false
}

// SetOsVersion gets a reference to the given string and assigns it to the OsVersion field.
func (o *Package) SetOsVersion(v string) {
	o.OsVersion = &v
}

// GetRelease returns the Release field value if set, zero value otherwise.
func (o *Package) GetRelease() string {
	if o == nil || o.Release == nil {
		var ret string
		return ret
	}
	return *o.Release
}

// GetReleaseOk returns a tuple with the Release field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Package) GetReleaseOk() (*string, bool) {
	if o == nil || o.Release == nil {
		return nil, false
	}
	return o.Release, true
}

// HasRelease returns a boolean if a field has been set.
func (o *Package) HasRelease() bool {
	if o != nil && o.Release != nil {
		return true
	}

	return false
}

// SetRelease gets a reference to the given string and assigns it to the Release field.
func (o *Package) SetRelease(v string) {
	o.Release = &v
}

// GetSize returns the Size field value if set, zero value otherwise.
func (o *Package) GetSize() int64 {
	if o == nil || o.Size == nil {
		var ret int64
		return ret
	}
	return *o.Size
}

// GetSizeOk returns a tuple with the Size field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Package) GetSizeOk() (*int64, bool) {
	if o == nil || o.Size == nil {
		return nil, false
	}
	return o.Size, true
}

// HasSize returns a boolean if a field has been set.
func (o *Package) HasSize() bool {
	if o != nil && o.Size != nil {
		return true
	}

	return false
}

// SetSize gets a reference to the given int64 and assigns it to the Size field.
func (o *Package) SetSize(v int64) {
	o.Size = &v
}

// GetSource returns the Source field value if set, zero value otherwise.
func (o *Package) GetSource() string {
	if o == nil || o.Source == nil {
		var ret string
		return ret
	}
	return *o.Source
}

// GetSourceOk returns a tuple with the Source field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Package) GetSourceOk() (*string, bool) {
	if o == nil || o.Source == nil {
		return nil, false
	}
	return o.Source, true
}

// HasSource returns a boolean if a field has been set.
func (o *Package) HasSource() bool {
	if o != nil && o.Source != nil {
		return true
	}

	return false
}

// SetSource gets a reference to the given string and assigns it to the Source field.
func (o *Package) SetSource(v string) {
	o.Source = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *Package) GetType() string {
	if o == nil || o.Type == nil {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Package) GetTypeOk() (*string, bool) {
	if o == nil || o.Type == nil {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *Package) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *Package) SetType(v string) {
	o.Type = &v
}

// GetVersion returns the Version field value if set, zero value otherwise.
func (o *Package) GetVersion() string {
	if o == nil || o.Version == nil {
		var ret string
		return ret
	}
	return *o.Version
}

// GetVersionOk returns a tuple with the Version field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Package) GetVersionOk() (*string, bool) {
	if o == nil || o.Version == nil {
		return nil, false
	}
	return o.Version, true
}

// HasVersion returns a boolean if a field has been set.
func (o *Package) HasVersion() bool {
	if o != nil && o.Version != nil {
		return true
	}

	return false
}

// SetVersion gets a reference to the given string and assigns it to the Version field.
func (o *Package) SetVersion(v string) {
	o.Version = &v
}

func (o Package) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Assessment != nil {
		toSerialize["assessment"] = o.Assessment
	}
	if o.Description != nil {
		toSerialize["description"] = o.Description
	}
	if o.Epoch != nil {
		toSerialize["epoch"] = o.Epoch
	}
	if o.HomePage != nil {
		toSerialize["home_page"] = o.HomePage
	}
	if o.Id != nil {
		toSerialize["id"] = o.Id
	}
	if o.License != nil {
		toSerialize["license"] = o.License
	}
	if o.Maintainer != nil {
		toSerialize["maintainer"] = o.Maintainer
	}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	if o.OsArchitecture != nil {
		toSerialize["os_architecture"] = o.OsArchitecture
	}
	if o.OsFamily != nil {
		toSerialize["os_family"] = o.OsFamily
	}
	if o.OsName != nil {
		toSerialize["os_name"] = o.OsName
	}
	if o.OsVendor != nil {
		toSerialize["os_vendor"] = o.OsVendor
	}
	if o.OsVersion != nil {
		toSerialize["os_version"] = o.OsVersion
	}
	if o.Release != nil {
		toSerialize["release"] = o.Release
	}
	if o.Size != nil {
		toSerialize["size"] = o.Size
	}
	if o.Source != nil {
		toSerialize["source"] = o.Source
	}
	if o.Type != nil {
		toSerialize["type"] = o.Type
	}
	if o.Version != nil {
		toSerialize["version"] = o.Version
	}
	return json.Marshal(toSerialize)
}

type NullablePackage struct {
	value *Package
	isSet bool
}

func (v NullablePackage) Get() *Package {
	return v.value
}

func (v *NullablePackage) Set(val *Package) {
	v.value = val
	v.isSet = true
}

func (v NullablePackage) IsSet() bool {
	return v.isSet
}

func (v *NullablePackage) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePackage(val *Package) *NullablePackage {
	return &NullablePackage{value: val, isSet: true}
}

func (v NullablePackage) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePackage) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
