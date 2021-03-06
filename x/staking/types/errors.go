// nolint
package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidValidator  CodeType = 101
	CodeInvalidDelegation CodeType = 102
	CodeInvalidInput      CodeType = 103
	CodeInvalidAddress             = sdk.CodeInvalidAddress
	CodeUnknownRequest             = sdk.CodeUnknownRequest

	CodeInvalidMinSelfDelegation CodeType = 104
	CodeInvalidProxy             CodeType = 105
	CodeInvalidShareAdding       CodeType = 106
)

// ErrNilValidatorAddr returns an error when an empty validator address appears
func ErrNilValidatorAddr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "validator address is nil")
}

// ErrBadValidatorAddr returns an error when an invalid validator address appears
func ErrBadValidatorAddr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAddress, "validator address is invalid")
}

// ErrNoValidatorFound returns an error when a validator doesn't exist
func ErrNoValidatorFound(codespace sdk.CodespaceType, valAddr string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "validator %s does not exist", valAddr)
}

// ErrValidatorOwnerExists returns an error when the validator address has been registered
func ErrValidatorOwnerExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator,
		"validator already exist for this operator address, must use new validator operator address")
}

// ErrValidatorPubKeyExists returns an error when the validator consensus pubkey has been registered
func ErrValidatorPubKeyExists(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator,
		"validator already exist for this pubkey, must use new validator pubkey")
}

// ErrValidatorPubKeyTypeNotSupported returns an error when the type of  pubkey was not supported
func ErrValidatorPubKeyTypeNotSupported(codespace sdk.CodespaceType, keyType string,
	supportedTypes []string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator,
		"validator pubkey type %s is not supported, must use %s", keyType, strings.Join(supportedTypes, ","))
}

// ErrDescriptionLength returns an error when the description of validator has a wrong length
func ErrDescriptionLength(codespace sdk.CodespaceType, descriptor string, got, max int) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator,
		"bad description length for %v, got length %v, max is %v", descriptor, got, max)
}

// ErrCommissionNegative returns an error when the commission is not positive
func ErrCommissionNegative(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "commission must be positive")
}

// ErrCommissionHuge returns an error when the commission is greater than 100%
func ErrCommissionHuge(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "commission cannot be more than 100%")
}

// ErrCommissionGTMaxRate returns an error when the commission rate is greater than the max rate
func ErrCommissionGTMaxRate(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "commission cannot be more than the max rate")
}

// ErrCommissionUpdateTime returns an error when the commission is remodified within 24 hours
func ErrCommissionUpdateTime(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "commission cannot be changed more than once in 24h")
}

// ErrCommissionChangeRateNegative returns an error when the commission change rate is not positive
func ErrCommissionChangeRateNegative(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "commission change rate must be positive")
}

// ErrCommissionChangeRateGTMaxRate returns an error when the commission change rate is greater than the max rate
func ErrCommissionChangeRateGTMaxRate(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "commission change rate cannot be more than the max rate")
}

// ErrCommissionGTMaxChangeRate returns an error when the new rate % points change is greater than the max change rate
func ErrCommissionGTMaxChangeRate(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "commission cannot be changed more than max change rate")
}

// ErrMinSelfDelegationInvalid returns an error when the msd isn't positive
func ErrMinSelfDelegationInvalid(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidValidator, "minimum self delegation must be a positive integer")
}

// ErrNilDelegatorAddr returns an error when the delegator address is nil
func ErrNilDelegatorAddr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "delegator address is nil")
}

// ErrWrongOperationAddr returns an error when the address is not expected
func ErrWrongOperationAddr(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, "wrong operation address found. msg: %s", msg)
}

// ErrBadDenom returns an error when the coin denomination is invalid
func ErrBadDenom(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDelegation, "invalid coin denomination")
}

// ErrBadDelegationAmount returns an error when the amount of delegation isn't positive
func ErrBadDelegationAmount(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDelegation, "amount must be > 0")
}

// TODO:remove it later
// ErrInsufficientShares returns an error when the delegation shares is not enough
func ErrInsufficientShares(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDelegation, "insufficient delegation shares")
}

// ErrNoUnbondingDelegation returns an error when the unbonding delegation doesn't exist
func ErrNoUnbondingDelegation(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDelegation, "no unbonding delegation found")
}

// ErrAddSharesToDismission returns an error when a zero-msd validator becomes the shares adding target
func ErrAddSharesToDismission(codespace sdk.CodespaceType, valAddr string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidShareAdding,
		"failed. destroyed validator %s isn't allowed to add shares to. please get rid of it from the "+
			"shares adding list by adding shares to other validators again or unbond all delegated tokens", valAddr)
}

// ErrNoAvailableValsToAddShares returns an error when none of the validators in shares adding list is available
func ErrNoAvailableValsToAddShares(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidShareAdding,
		"failed. there's no available validators among the shares adding list")
}

// ErrDelegatorNotAProxy returns an error when the target delegator to bind is not registered as a proxy
func ErrDelegatorNotAProxy(codespace sdk.CodespaceType, delegator string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDelegation,
		"failed. target address haven't reg as a proxy %s", delegator)
}

// ErrNeverProxied returns an error when a delegator who's not a proxy trys to unreg
func ErrNeverProxied(codespace sdk.CodespaceType, delegator string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidProxy,
		"failed. delegator %s has never registered as a proxy", delegator)
}

// ErrAlreadyProxied returns an error when a proxy trys to reg the second time
func ErrAlreadyProxied(codespace sdk.CodespaceType, delegator string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidProxy,
		"failed. delegator %s has already registered as a proxy", delegator)
}

// ErrAddSharesDuringProxy returns an error when a delegator who has bound tries to add shares to validators by itself
func ErrAddSharesDuringProxy(codespace sdk.CodespaceType, delegator string, proxy string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDelegation,
		"failed. ban from adding shares to validators before unbinding proxy relationship between %s and %s", delegator, proxy)
}

// ErrDoubleProxy returns an error when a delegator trys to bind more than one proxy
func ErrDoubleProxy(codespace sdk.CodespaceType, delegator string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDelegation,
		"failed. proxy isn't allowed to bind with other proxy %s", delegator)
}

// ErrNotFoundProxy returns an error when the proxy doesn't exist
func ErrNotFoundProxy(codespace sdk.CodespaceType, delegator string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDelegation,
		"failed. no proxy with %s", delegator)
}

// ErrInvalidDelegation returns an error when the delegation is invalid
func ErrInvalidDelegation(codespace sdk.CodespaceType, delegator string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDelegation,
		"failed. invalid delegation on %s", delegator)
}

// ErrNilValidatorAddrs returns an error when there's no target validator address in MsgAddShares
func ErrNilValidatorAddrs(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput,
		"failed. validator addresses are nil")
}

// ErrExceedValidatorAddrs returns an error when the number of target validators exceeds the max limit
func ErrExceedValidatorAddrs(codespace sdk.CodespaceType, num int) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput,
		"failed. the number of validator addresses is over the limit %d", num)
}

// ErrNoDelegationToAddShares returns an error when there's no delegation to support adding shares to validators
func ErrNoDelegationToAddShares(codespace sdk.CodespaceType, delegator string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDelegation,
		"failed. there's no delegation of %s", delegator)
}

// ErrNotInDelegating returns an error when the UndelegationInfo doesn't exist during it's unbonding period
func ErrNotInDelegating(codespace sdk.CodespaceType, addr string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDelegation,
		"failed. the addr %s is not in the status of undelegating", addr)
}

// ErrInsufficientDelegation returns an error when the delegation left is not enough for unbonding
func ErrInsufficientDelegation(codespace sdk.CodespaceType, quantity, delLeft string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDelegation,
		"failed. insufficient delegation. [delegation left]:%s, [quantity to unbond]:%s", delLeft, quantity)
}

// ErrInsufficientQuantity returns an error when the quantity is less than the min delegation limit
func ErrInsufficientQuantity(codespace sdk.CodespaceType, quantity, minLimit string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDelegation,
		"failed. insufficient quantity. [min limit]:%s, [quantity]:%s", minLimit, quantity)
}

// ErrMoreMinSelfDelegation returns an error when the msd doesn't match the rest of shares on a validator
func ErrMoreMinSelfDelegation(codespace sdk.CodespaceType, valAddr string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidMinSelfDelegation,
		"failed. min self delegation of %s is more than its shares", valAddr)
}

// ErrNoMinSelfDelegation returns an error when the msd has already been unbonded
func ErrNoMinSelfDelegation(codespace sdk.CodespaceType, valAddr string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidMinSelfDelegation,
		"failed. there's no min self delegation on %s", valAddr)
}

// ErrBadUnDelegationAmount returns an error when the amount of delegation is not positive
func ErrBadUnDelegationAmount(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDelegation,
		"failed. amount must be greater than 0")
}

// ErrInvalidProxyUpdating returns an error when the total delegated tokens on a proxy are going to be negative
func ErrInvalidProxyUpdating(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidProxy,
		"failed. the total delegated tokens on the proxy will be negative after this update")
}

// ErrInvalidProxyWithdrawTotal returns an error when proxy withdraws total tokens
func ErrInvalidProxyWithdrawTotal(codespace sdk.CodespaceType, addr string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidProxy,
		"failed. proxy %s has to unreg before withdrawing total tokens", addr)
}

// ErrAlreadyAddedShares returns an error when a delegator tries to bind proxy after adding shares
func ErrAlreadyAddedShares(codespace sdk.CodespaceType, delAddr string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidProxy,
		"failed. delegator %s isn't allowed to bind proxy while it has added shares. please unbond the delegation first", delAddr)
}

// ErrNoDelegatorExisted returns an error when the info if a certain delegator doesn't exist
func ErrNoDelegatorExisted(codespace sdk.CodespaceType, delAddr string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDelegation,
		"failed. delegator %s doesn't exist", delAddr)
}

// ErrTargetValsDuplicate returns an error when the target validators in voting list are duplicate
func ErrTargetValsDuplicate(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidShareAdding,
		"failed. duplicate target validators")
}

// ErrAlreadyBound returns an error when a delegator keeps binding a proxy before proxy register
func ErrAlreadyBound(codespace sdk.CodespaceType, delAddr string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidProxy,
		"failed. %s has already bound a proxy. it's necessary to unbind before proxy register",
		delAddr)
}
