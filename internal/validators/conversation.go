package validators

import (
    "errors"
    "strconv"
    "strings"

    "github.com/lincentpega/pcrm/internal/dto"
)

func ValidateConversationID(idStr string) (int64, error) {
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        return 0, errors.New("invalid conversation ID")
    }
    return id, nil
}

func ValidateConversationRequest(req *dto.ConversationRequest) error {
    if req.ConversationTypeID <= 0 {
        return errors.New("conversation type id is required")
    }
    if req.Initiator == "" {
        return errors.New("initiator is required")
    }
    v := strings.ToLower(req.Initiator)
    if v != "owner" && v != "person" {
        return errors.New("initiator must be 'owner' or 'person'")
    }
    if strings.TrimSpace(req.Notes) == "" {
        return errors.New("notes is required")
    }
    req.Initiator = v
    return nil
}

