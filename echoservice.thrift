namespace java com.motus.echoservice.thrift

/**
 * The EchoService Error Code. <br/>
 */
enum EchoServiceErrorCode {
	ERROR_CODE_0,
	ERROR_CODE_1,
	// TODO Add your desired error codes here
	OTHER_ERROR_REASON,
	
	/**
	 * Unknown error (unexpected ones)
	 */
	UNKNOWN_ERROR;
}

/**
 * The EchoService Business Exception. <br/>
 * We send this back in case there is an error in actual business logic. <br/>
 * This is supposed to be non-retry-able from client side.
 */
exception TEchoServiceBusinessException
{
	1: optional map<EchoServiceErrorCode, string> errors, // would be nice to say something about what happened
	2: optional string serviceMessage,
	3: optional string serviceStackTrace;
}

/**
 * The EchoService Failure Exception. <br/>
 * We send this back in case there is an error in communication with other services and so on. <br/>
 * This is supposed to be retry-able from client side.
 */
exception TEchoServiceFailureException
{
	1: optional map<EchoServiceErrorCode, string> errors, // would be nice to say something about what happened
	2: optional string serviceMessage,
	3: optional string serviceStackTrace;
}

/**
 * Echo service input DTO. <br/>
 */
struct TEchoServiceInputDTO
{
	1: string message;
}

/**
 * Echo service output DTO.<br/>
 */
struct TEchoServiceOutputDTO
{
	1: string echoMessage;
}

struct EchoEvent
{
    1: string echoMessage
}

/**
 * Thrift service interface for the EchoService.
 */
service EchoService {

    /**
    * Sample sayMyName method.<br/>
    * It shows how to handle input and output as DTOs.
    */
    TEchoServiceOutputDTO echo(1: TEchoServiceInputDTO inputDTO) throws (1: TEchoServiceBusinessException businessException, 2: TEchoServiceFailureException failureException);
}