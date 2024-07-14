package controller

import (
	"github.com/google/uuid"

	"chiko/pkg/entity"
)

// DeleteBookmark is used to delete a bookmark from the bookmark tree and save the bookmark
func (c Controller) DeleteBookmark(b entity.Session) error {
	// Helper function to remove a session from a bookmark
	removeSession := func(sessions []entity.Session, sessionID *uuid.UUID) []entity.Session {
		for i, session := range sessions {
			if session.ID == *sessionID {
				return append(sessions[:i], sessions[i+1:]...)
			}
		}
		return sessions
	}

	// Regenerate the boomark tree
	for i, bookmark := range *c.Bookmarks {
		updatedSessions := removeSession(bookmark.Sessions, &b.ID)
		if len(updatedSessions) != len(bookmark.Sessions) {
			bookmark.Sessions = updatedSessions
			if len(bookmark.Sessions) == 0 {
				*c.Bookmarks = append((*c.Bookmarks)[:i], (*c.Bookmarks)[i+1:]...)
			}
			err := c.SaveBookmark()
			if err != nil {
				return err
			}

			return nil
		}
	}

	return nil
}
