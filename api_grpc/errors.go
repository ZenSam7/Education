package api_grpc

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

func unauthenticatedError(err error) error {
	return status.Errorf(codes.Unauthenticated, "пользователь не авторизовался: %s", err)
}

func fieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{Field: field, Description: err.Error()}
}
