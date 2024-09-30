// These are not tied to a particular semver and are intended to be used in packages which otherwise have no
// dependency nor need for dependency upon otel libraries.
// Taken from https://github.com/open-telemetry/opentelemetry-collector/blob/52e11ecb8293dc9ae75edbabb4f0c176b386d9c7/semconv/v1.27.0/generated_attribute_group.go

package errors

const (
	OtelClientAddress                   = "client.address"
	OtelClientPort                      = "client.port"
	OtelCodeFilePath                    = "code.filepath"
	OtelCodeFunction                    = "code.function"
	OtelCodeLineNo                      = "code.lineno"
	OtelCodeNamespace                   = "code.namespace"
	OtelFileDirectory                   = "file.directory"
	OtelFileExtension                   = "file.extension"
	OtelFileName                        = "file.name"
	OtelFilePath                        = "file.path"
	OtelFileSize                        = "file.size"
	OtelHostID                          = "host.id"
	OtelHostIP                          = "host.ip"
	OtelHTTPRequestBodySize             = "http.request.body.size"
	OtelHTTPRequestMethod               = "http.request.method"
	OtelHTTPRequestSize                 = "http.request.size"
	OtelHTTPResponseBodySize            = "http.response.body.size"
	OtelHTTPResponseSize                = "http.response.size"
	OtelHTTPResponseStatusCode          = "http.response.status_code"
	OtelHTTPUserAgentName               = "user_agent.name"
	OtelURLDomain                       = "url.domain" // 'www.foo.bar', 'opentelemetry.io', '3.12.167.2',
	OtelURLFull                         = "url.full"   // https://www.foo.bar/search?q=OpenTelemetry#SemConv
	OtelURLPath                         = "url.path"   // /search /query
	OtelURLPort                         = "url.port"   // 80 443
	OtelURLQuery                        = "url.query"  // q=OpenTelemetry your=mom
	OtelURLScheme                       = "url.scheme" // 'https', 'http'
	OtelMessagingClientID               = "messaging.client.id"
	OtelMessagingConsumerGroupName      = "messaging.consumer.group.name"
	OtelMessagingDestinationPartitionID = "messaging.destination.partition.id"
	OtelMessagingMessageBodySize        = "messaging.message.body.size"
	OtelMessagingMessageConversationID  = "messaging.message.conversation_id"
	OtelMessagingMessageEnvelopeSize    = "messaging.message.envelope.size"
	OtelMessagingMessageID              = "messaging.message.id"
	OtelMessagingOperationName          = "messaging.operation.name" // ack, nack, send, etc..
	OtelMessagingOperationType          = "messaging.operation.type" // publish, create, receive, settle, deliver
	OtelMessagingSystem                 = "messaging.system"
	OtelNetworkConnectionType           = "network.connection.type"
	OtelNetworkLocalAddress             = "network.local.address"
	OtelNetworkLocalPort                = "network.local.port"
	OtelNetworkPeerAddress              = "network.peer.address"
	OtelNetworkPeerPort                 = "network.peer.port"
	OtelNetworkProtocolName             = "network.protocol.name" // http, amqp, mqtt, etc..
	OtelNetworkTransport                = "network.transport"     // tcp, udp
	OtelNetworkType                     = "network.type"          // ipv4, ipv6
	OtelServerAddress                   = "server.address"
	OtelServerPort                      = "server.port"
	OtelServiceInstanceID               = "service.instance.id"
	OtelServiceName                     = "service.name"
	OtelServiceNamespace                = "service.namespace"
	OtelServiceVersion                  = "service.version"
	OtelSessionID                       = "session.id"
	OtelTLSCipher                       = "tls.cipher"
	OtelTLSProtocolVersion              = "tls.protocol.version"
	OtelTLSServerSubject                = "tls.server.subject"
	OtelUserEmail                       = "user.email"
	OtelUserID                          = "user.id"
	OtelUserName                        = "user.name"
	OtelUserRoles                       = "user.roles"
)
