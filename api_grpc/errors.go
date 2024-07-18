package api_grpc

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func fieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{Field: field, Description: err.Error()}
}

func wrapFeildErrors(violation []*errdetails.BadRequest_FieldViolation) error {
	if len(violation) != 0 {
		statusInvalid := status.New(codes.InvalidArgument, "неправильные параметры")

		statusDetails, err := statusInvalid.WithDetails(&errdetails.BadRequest{FieldViolations: violation})
		if err != nil {
			return statusInvalid.Err()
		}

		return statusDetails.Err()
	}

	return nil
}
