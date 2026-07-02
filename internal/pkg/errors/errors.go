package apperrors

import "errors"

var (
	ErrNotFound           = errors.New("not_found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidOTP         = errors.New("invalid_otp")
	ErrOTPExpired         = errors.New("otp_expired")
	ErrInvalidPhone       = errors.New("invalid_phone")
	ErrUserExists         = errors.New("user_exists")
	ErrUserNotFound       = errors.New("user_not_found")
	ErrEventNotFound      = errors.New("event_not_found")
	ErrGiftNotFound       = errors.New("gift_not_found")
	ErrGiftReserved       = errors.New("gift_reserved")
	ErrInsufficientAmount = errors.New("insufficient_amount")
	ErrValidation         = errors.New("validation_error")
	ErrEventNotOpen       = errors.New("event_not_open")
	ErrEventNotWithdrawable = errors.New("event_not_withdrawable")
	ErrInvalidWithdrawal  = errors.New("invalid_withdrawal")
	ErrInvalidAdminKey    = errors.New("invalid_admin_key")
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	MessageAR string `json:"message_ar,omitempty"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code, message, messageAR string, err error) *AppError {
	return &AppError{Code: code, Message: message, MessageAR: messageAR, Err: err}
}

func Map(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	switch {
	case errors.Is(err, ErrNotFound), errors.Is(err, ErrUserNotFound), errors.Is(err, ErrEventNotFound), errors.Is(err, ErrGiftNotFound):
		return New("NOT_FOUND", "Resource not found", "المورد غير موجود", err)
	case errors.Is(err, ErrUnauthorized):
		return New("UNAUTHORIZED", "Unauthorized", "غير مصرح", err)
	case errors.Is(err, ErrForbidden):
		return New("FORBIDDEN", "Forbidden", "محظور", err)
	case errors.Is(err, ErrInvalidOTP):
		return New("INVALID_OTP", "Invalid verification code", "رمز التحقق غير صحيح", err)
	case errors.Is(err, ErrOTPExpired):
		return New("OTP_EXPIRED", "Verification code expired", "انتهت صلاحية رمز التحقق", err)
	case errors.Is(err, ErrInvalidPhone):
		return New("INVALID_PHONE", "Invalid Saudi phone number (+966)", "رقم الجوال غير صالح (+966)", err)
	case errors.Is(err, ErrUserExists):
		return New("USER_EXISTS", "User already registered", "المستخدم مسجل مسبقاً", err)
	case errors.Is(err, ErrGiftReserved):
		return New("GIFT_RESERVED", "Gift is already reserved", "الهدية محجوزة مسبقاً", err)
	case errors.Is(err, ErrInsufficientAmount):
		return New("INSUFFICIENT_AMOUNT", "Contribution exceeds remaining amount", "المبلغ يتجاوز المتبقي", err)
	case errors.Is(err, ErrEventNotOpen):
		return New("EVENT_NOT_OPEN", "Event is not open for contributions", "المناسبة غير مفتوحة للمساهمات", err)
	case errors.Is(err, ErrEventNotWithdrawable):
		return New("EVENT_NOT_WITHDRAWABLE", "Event is not ready for withdrawal", "المناسبة غير جاهزة للسحب", err)
	case errors.Is(err, ErrInvalidWithdrawal):
		return New("INVALID_WITHDRAWAL", "Invalid withdrawal amount", "مبلغ السحب غير صالح", err)
	case errors.Is(err, ErrInvalidAdminKey):
		return New("UNAUTHORIZED", "Invalid admin credentials", "بيانات المسؤول غير صحيحة", err)
	case errors.Is(err, ErrValidation):
		return New("VALIDATION_ERROR", err.Error(), err.Error(), err)
	default:
		return New("INTERNAL_ERROR", "Internal server error", "خطأ في الخادم", err)
	}
}
