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

// PackageCvssV3 ${package.cvss.v3.description}, ${packagecvssv3.description}
type PackageCvssV3 struct {
	// ${package.cvss.v3.attack_complexity.description}, ${package.cvss.v3.attack.complexity.description}, ${packagecvssv3.attack_complexity.description}, ${packagecvssv3.attack.complexity.description}
	AttackComplexity *string `json:"attack_complexity,omitempty"`
	// ${packagecvssv3.attack_vector.description}, ${package.cvss.v3.attack_vector.description}, ${packagecvssv3.attack.vector.description}, ${package.cvss.v3.attack.vector.description}
	AttackVector *string `json:"attack_vector,omitempty"`
	// ${package.cvss.v3.availability_impact.description}, ${packagecvssv3.availability.impact.description}, ${packagecvssv3.availability_impact.description}, ${package.cvss.v3.availability.impact.description}
	AvailabilityImpact *string `json:"availability_impact,omitempty"`
	// ${packagecvssv3.confidentiality_impact.description}, ${packagecvssv3.confidentiality.impact.description}, ${package.cvss.v3.confidentiality_impact.description}, ${package.cvss.v3.confidentiality.impact.description}
	ConfidentialityImpact *string `json:"confidentiality_impact,omitempty"`
	// ${packagecvssv3.exploit_score.description}, ${packagecvssv3.exploit.score.description}, ${package.cvss.v3.exploit.score.description}, ${package.cvss.v3.exploit_score.description}
	ExploitScore *float64 `json:"exploit_score,omitempty"`
	// ${package.cvss.v3.impact.score.description}, ${packagecvssv3.impact_score.description}, ${package.cvss.v3.impact_score.description}, ${packagecvssv3.impact.score.description}
	ImpactScore *float64 `json:"impact_score,omitempty"`
	// ${packagecvssv3.integrity_impact.description}, ${package.cvss.v3.integrity.impact.description}, ${package.cvss.v3.integrity_impact.description}, ${packagecvssv3.integrity.impact.description}
	IntegrityImpact *string `json:"integrity_impact,omitempty"`
	// ${package.cvss.v3.privileges.required.description}, ${packagecvssv3.privileges_required.description}, ${package.cvss.v3.privileges_required.description}, ${packagecvssv3.privileges.required.description}
	PrivilegesRequired *string `json:"privileges_required,omitempty"`
	// ${package.cvss.v3.scope.description}, ${packagecvssv3.scope.description}
	Scope *string `json:"scope,omitempty"`
	// ${packagecvssv3.score.description}, ${package.cvss.v3.score.description}
	Score *float64 `json:"score,omitempty"`
	// ${package.cvss.v3.user.interaction.description}, ${package.cvss.v3.user_interaction.description}, ${packagecvssv3.user_interaction.description}, ${packagecvssv3.user.interaction.description}
	UserInteraction *string `json:"user_interaction,omitempty"`
	// ${packagecvssv3.vector.description}, ${package.cvss.v3.vector.description}
	Vector *string `json:"vector,omitempty"`
}

// NewPackageCvssV3 instantiates a new PackageCvssV3 object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPackageCvssV3() *PackageCvssV3 {
	this := PackageCvssV3{}
	return &this
}

// NewPackageCvssV3WithDefaults instantiates a new PackageCvssV3 object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPackageCvssV3WithDefaults() *PackageCvssV3 {
	this := PackageCvssV3{}
	return &this
}

// GetAttackComplexity returns the AttackComplexity field value if set, zero value otherwise.
func (o *PackageCvssV3) GetAttackComplexity() string {
	if o == nil || o.AttackComplexity == nil {
		var ret string
		return ret
	}
	return *o.AttackComplexity
}

// GetAttackComplexityOk returns a tuple with the AttackComplexity field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PackageCvssV3) GetAttackComplexityOk() (*string, bool) {
	if o == nil || o.AttackComplexity == nil {
		return nil, false
	}
	return o.AttackComplexity, true
}

// HasAttackComplexity returns a boolean if a field has been set.
func (o *PackageCvssV3) HasAttackComplexity() bool {
	if o != nil && o.AttackComplexity != nil {
		return true
	}

	return false
}

// SetAttackComplexity gets a reference to the given string and assigns it to the AttackComplexity field.
func (o *PackageCvssV3) SetAttackComplexity(v string) {
	o.AttackComplexity = &v
}

// GetAttackVector returns the AttackVector field value if set, zero value otherwise.
func (o *PackageCvssV3) GetAttackVector() string {
	if o == nil || o.AttackVector == nil {
		var ret string
		return ret
	}
	return *o.AttackVector
}

// GetAttackVectorOk returns a tuple with the AttackVector field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PackageCvssV3) GetAttackVectorOk() (*string, bool) {
	if o == nil || o.AttackVector == nil {
		return nil, false
	}
	return o.AttackVector, true
}

// HasAttackVector returns a boolean if a field has been set.
func (o *PackageCvssV3) HasAttackVector() bool {
	if o != nil && o.AttackVector != nil {
		return true
	}

	return false
}

// SetAttackVector gets a reference to the given string and assigns it to the AttackVector field.
func (o *PackageCvssV3) SetAttackVector(v string) {
	o.AttackVector = &v
}

// GetAvailabilityImpact returns the AvailabilityImpact field value if set, zero value otherwise.
func (o *PackageCvssV3) GetAvailabilityImpact() string {
	if o == nil || o.AvailabilityImpact == nil {
		var ret string
		return ret
	}
	return *o.AvailabilityImpact
}

// GetAvailabilityImpactOk returns a tuple with the AvailabilityImpact field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PackageCvssV3) GetAvailabilityImpactOk() (*string, bool) {
	if o == nil || o.AvailabilityImpact == nil {
		return nil, false
	}
	return o.AvailabilityImpact, true
}

// HasAvailabilityImpact returns a boolean if a field has been set.
func (o *PackageCvssV3) HasAvailabilityImpact() bool {
	if o != nil && o.AvailabilityImpact != nil {
		return true
	}

	return false
}

// SetAvailabilityImpact gets a reference to the given string and assigns it to the AvailabilityImpact field.
func (o *PackageCvssV3) SetAvailabilityImpact(v string) {
	o.AvailabilityImpact = &v
}

// GetConfidentialityImpact returns the ConfidentialityImpact field value if set, zero value otherwise.
func (o *PackageCvssV3) GetConfidentialityImpact() string {
	if o == nil || o.ConfidentialityImpact == nil {
		var ret string
		return ret
	}
	return *o.ConfidentialityImpact
}

// GetConfidentialityImpactOk returns a tuple with the ConfidentialityImpact field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PackageCvssV3) GetConfidentialityImpactOk() (*string, bool) {
	if o == nil || o.ConfidentialityImpact == nil {
		return nil, false
	}
	return o.ConfidentialityImpact, true
}

// HasConfidentialityImpact returns a boolean if a field has been set.
func (o *PackageCvssV3) HasConfidentialityImpact() bool {
	if o != nil && o.ConfidentialityImpact != nil {
		return true
	}

	return false
}

// SetConfidentialityImpact gets a reference to the given string and assigns it to the ConfidentialityImpact field.
func (o *PackageCvssV3) SetConfidentialityImpact(v string) {
	o.ConfidentialityImpact = &v
}

// GetExploitScore returns the ExploitScore field value if set, zero value otherwise.
func (o *PackageCvssV3) GetExploitScore() float64 {
	if o == nil || o.ExploitScore == nil {
		var ret float64
		return ret
	}
	return *o.ExploitScore
}

// GetExploitScoreOk returns a tuple with the ExploitScore field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PackageCvssV3) GetExploitScoreOk() (*float64, bool) {
	if o == nil || o.ExploitScore == nil {
		return nil, false
	}
	return o.ExploitScore, true
}

// HasExploitScore returns a boolean if a field has been set.
func (o *PackageCvssV3) HasExploitScore() bool {
	if o != nil && o.ExploitScore != nil {
		return true
	}

	return false
}

// SetExploitScore gets a reference to the given float64 and assigns it to the ExploitScore field.
func (o *PackageCvssV3) SetExploitScore(v float64) {
	o.ExploitScore = &v
}

// GetImpactScore returns the ImpactScore field value if set, zero value otherwise.
func (o *PackageCvssV3) GetImpactScore() float64 {
	if o == nil || o.ImpactScore == nil {
		var ret float64
		return ret
	}
	return *o.ImpactScore
}

// GetImpactScoreOk returns a tuple with the ImpactScore field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PackageCvssV3) GetImpactScoreOk() (*float64, bool) {
	if o == nil || o.ImpactScore == nil {
		return nil, false
	}
	return o.ImpactScore, true
}

// HasImpactScore returns a boolean if a field has been set.
func (o *PackageCvssV3) HasImpactScore() bool {
	if o != nil && o.ImpactScore != nil {
		return true
	}

	return false
}

// SetImpactScore gets a reference to the given float64 and assigns it to the ImpactScore field.
func (o *PackageCvssV3) SetImpactScore(v float64) {
	o.ImpactScore = &v
}

// GetIntegrityImpact returns the IntegrityImpact field value if set, zero value otherwise.
func (o *PackageCvssV3) GetIntegrityImpact() string {
	if o == nil || o.IntegrityImpact == nil {
		var ret string
		return ret
	}
	return *o.IntegrityImpact
}

// GetIntegrityImpactOk returns a tuple with the IntegrityImpact field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PackageCvssV3) GetIntegrityImpactOk() (*string, bool) {
	if o == nil || o.IntegrityImpact == nil {
		return nil, false
	}
	return o.IntegrityImpact, true
}

// HasIntegrityImpact returns a boolean if a field has been set.
func (o *PackageCvssV3) HasIntegrityImpact() bool {
	if o != nil && o.IntegrityImpact != nil {
		return true
	}

	return false
}

// SetIntegrityImpact gets a reference to the given string and assigns it to the IntegrityImpact field.
func (o *PackageCvssV3) SetIntegrityImpact(v string) {
	o.IntegrityImpact = &v
}

// GetPrivilegesRequired returns the PrivilegesRequired field value if set, zero value otherwise.
func (o *PackageCvssV3) GetPrivilegesRequired() string {
	if o == nil || o.PrivilegesRequired == nil {
		var ret string
		return ret
	}
	return *o.PrivilegesRequired
}

// GetPrivilegesRequiredOk returns a tuple with the PrivilegesRequired field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PackageCvssV3) GetPrivilegesRequiredOk() (*string, bool) {
	if o == nil || o.PrivilegesRequired == nil {
		return nil, false
	}
	return o.PrivilegesRequired, true
}

// HasPrivilegesRequired returns a boolean if a field has been set.
func (o *PackageCvssV3) HasPrivilegesRequired() bool {
	if o != nil && o.PrivilegesRequired != nil {
		return true
	}

	return false
}

// SetPrivilegesRequired gets a reference to the given string and assigns it to the PrivilegesRequired field.
func (o *PackageCvssV3) SetPrivilegesRequired(v string) {
	o.PrivilegesRequired = &v
}

// GetScope returns the Scope field value if set, zero value otherwise.
func (o *PackageCvssV3) GetScope() string {
	if o == nil || o.Scope == nil {
		var ret string
		return ret
	}
	return *o.Scope
}

// GetScopeOk returns a tuple with the Scope field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PackageCvssV3) GetScopeOk() (*string, bool) {
	if o == nil || o.Scope == nil {
		return nil, false
	}
	return o.Scope, true
}

// HasScope returns a boolean if a field has been set.
func (o *PackageCvssV3) HasScope() bool {
	if o != nil && o.Scope != nil {
		return true
	}

	return false
}

// SetScope gets a reference to the given string and assigns it to the Scope field.
func (o *PackageCvssV3) SetScope(v string) {
	o.Scope = &v
}

// GetScore returns the Score field value if set, zero value otherwise.
func (o *PackageCvssV3) GetScore() float64 {
	if o == nil || o.Score == nil {
		var ret float64
		return ret
	}
	return *o.Score
}

// GetScoreOk returns a tuple with the Score field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PackageCvssV3) GetScoreOk() (*float64, bool) {
	if o == nil || o.Score == nil {
		return nil, false
	}
	return o.Score, true
}

// HasScore returns a boolean if a field has been set.
func (o *PackageCvssV3) HasScore() bool {
	if o != nil && o.Score != nil {
		return true
	}

	return false
}

// SetScore gets a reference to the given float64 and assigns it to the Score field.
func (o *PackageCvssV3) SetScore(v float64) {
	o.Score = &v
}

// GetUserInteraction returns the UserInteraction field value if set, zero value otherwise.
func (o *PackageCvssV3) GetUserInteraction() string {
	if o == nil || o.UserInteraction == nil {
		var ret string
		return ret
	}
	return *o.UserInteraction
}

// GetUserInteractionOk returns a tuple with the UserInteraction field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PackageCvssV3) GetUserInteractionOk() (*string, bool) {
	if o == nil || o.UserInteraction == nil {
		return nil, false
	}
	return o.UserInteraction, true
}

// HasUserInteraction returns a boolean if a field has been set.
func (o *PackageCvssV3) HasUserInteraction() bool {
	if o != nil && o.UserInteraction != nil {
		return true
	}

	return false
}

// SetUserInteraction gets a reference to the given string and assigns it to the UserInteraction field.
func (o *PackageCvssV3) SetUserInteraction(v string) {
	o.UserInteraction = &v
}

// GetVector returns the Vector field value if set, zero value otherwise.
func (o *PackageCvssV3) GetVector() string {
	if o == nil || o.Vector == nil {
		var ret string
		return ret
	}
	return *o.Vector
}

// GetVectorOk returns a tuple with the Vector field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PackageCvssV3) GetVectorOk() (*string, bool) {
	if o == nil || o.Vector == nil {
		return nil, false
	}
	return o.Vector, true
}

// HasVector returns a boolean if a field has been set.
func (o *PackageCvssV3) HasVector() bool {
	if o != nil && o.Vector != nil {
		return true
	}

	return false
}

// SetVector gets a reference to the given string and assigns it to the Vector field.
func (o *PackageCvssV3) SetVector(v string) {
	o.Vector = &v
}

func (o PackageCvssV3) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.AttackComplexity != nil {
		toSerialize["attack_complexity"] = o.AttackComplexity
	}
	if o.AttackVector != nil {
		toSerialize["attack_vector"] = o.AttackVector
	}
	if o.AvailabilityImpact != nil {
		toSerialize["availability_impact"] = o.AvailabilityImpact
	}
	if o.ConfidentialityImpact != nil {
		toSerialize["confidentiality_impact"] = o.ConfidentialityImpact
	}
	if o.ExploitScore != nil {
		toSerialize["exploit_score"] = o.ExploitScore
	}
	if o.ImpactScore != nil {
		toSerialize["impact_score"] = o.ImpactScore
	}
	if o.IntegrityImpact != nil {
		toSerialize["integrity_impact"] = o.IntegrityImpact
	}
	if o.PrivilegesRequired != nil {
		toSerialize["privileges_required"] = o.PrivilegesRequired
	}
	if o.Scope != nil {
		toSerialize["scope"] = o.Scope
	}
	if o.Score != nil {
		toSerialize["score"] = o.Score
	}
	if o.UserInteraction != nil {
		toSerialize["user_interaction"] = o.UserInteraction
	}
	if o.Vector != nil {
		toSerialize["vector"] = o.Vector
	}
	return json.Marshal(toSerialize)
}

type NullablePackageCvssV3 struct {
	value *PackageCvssV3
	isSet bool
}

func (v NullablePackageCvssV3) Get() *PackageCvssV3 {
	return v.value
}

func (v *NullablePackageCvssV3) Set(val *PackageCvssV3) {
	v.value = val
	v.isSet = true
}

func (v NullablePackageCvssV3) IsSet() bool {
	return v.isSet
}

func (v *NullablePackageCvssV3) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePackageCvssV3(val *PackageCvssV3) *NullablePackageCvssV3 {
	return &NullablePackageCvssV3{value: val, isSet: true}
}

func (v NullablePackageCvssV3) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePackageCvssV3) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
