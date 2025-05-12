package b2c

import (
	"errors"
	"gorm.io/gorm"
	repo "solution/internal/repository/b2c"
	"solution/internal/shared/models/b2c"
	"solution/internal/shared/models/b2c/dto"
	"strconv"
	"time"
)

type PromoService interface {
	GetPromosForUser(userID string, limit, offset int, category string, active *bool) ([]dto.PromoForUser, int64, error)
	GetPromo(promoID, userID string) (*dto.PromoForUser, error)
	LikePromo(promoID, userID string) error
	UnlikePromo(promoID, userID string) error
	AddComment(userID, promoID, text string) (*dto.CommentResponse, error)
	GetComments(promoID, limit, offset string) ([]dto.CommentResponse, int64, error)
	GetComment(promoID, commentID string) (*dto.CommentResponse, error)
	EditComment(userID, promoID, commentID, text string) (*dto.CommentResponse, error)
	DeleteComment(userID, promoID, commentID string) error
}

type promoService struct {
	repo repo.PromoRepository
}

func NewPromoService(repo repo.PromoRepository) PromoService {
	return &promoService{repo: repo}
}

func (s *promoService) GetPromosForUser(userID string, limit, offset int, category string, active *bool) ([]dto.PromoForUser, int64, error) {
	return s.repo.GetPromosForUser(userID, limit, offset, category, active)
}

func (s *promoService) GetPromo(promoID, userID string) (*dto.PromoForUser, error) {
	promo, err := s.repo.GetPromoForUserByID(promoID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrNotFound
		}
		return nil, err
	}

	return promo, nil
}

func (s *promoService) LikePromo(promoID, userID string) error {
	return s.repo.LikePromo(promoID, userID)
}

func (s *promoService) UnlikePromo(promoID, userID string) error {
	return s.repo.UnlikePromo(promoID, userID)
}

func (s *promoService) AddComment(userID, promoID, text string) (*dto.CommentResponse, error) {
	_, err := s.repo.GetPromoByID(promoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrNotFound
		}
		return nil, err
	}

	// Создаем новый комментарий
	comment := &b2c.Comment{
		UserID:  userID,
		PromoID: promoID,
		Text:    text,
	}

	err = s.repo.AddComment(comment)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	response := &dto.CommentResponse{
		ID:   comment.ID,
		Text: comment.Text,
		Date: comment.CreatedAt.Format(time.RFC3339),
		Author: dto.Author{
			Name:      user.Name,
			Surname:   user.Surname,
			AvatarURL: user.AvatarURL,
		},
	}

	return response, nil
}

// GetComments получает список комментариев к промокоду
func (s *promoService) GetComments(promoID, limitStr, offsetStr string) ([]dto.CommentResponse, int64, error) {
	_, err := s.repo.GetPromoByID(promoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, dto.ErrNotFound
		}
		return nil, 0, err
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	comments, totalCount, err := s.repo.GetComments(promoID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Формируем ответ
	var response []dto.CommentResponse
	for _, comment := range comments {
		user, err := s.repo.GetUserByID(comment.UserID)
		if err != nil {
			return nil, 0, err
		}

		commentResponse := dto.CommentResponse{
			ID:   comment.ID,
			Text: comment.Text,
			Date: comment.CreatedAt.Format(time.RFC3339),
			Author: dto.Author{
				Name:      user.Name,
				Surname:   user.Surname,
				AvatarURL: user.AvatarURL,
			},
		}

		response = append(response, commentResponse)
	}

	return response, totalCount, nil
}

// GetComment получает конкретный комментарий по его ID
func (s *promoService) GetComment(promoID, commentID string) (*dto.CommentResponse, error) {
	_, err := s.repo.GetPromoByID(promoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrNotFound
		}
		return nil, err
	}

	comment, err := s.repo.GetCommentByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrNotFound
		}
		return nil, err
	}

	// Проверяем, что комментарий относится к данному промокоду
	if comment.PromoID != promoID {
		return nil, dto.ErrNotFound
	}

	// Получаем данные пользователя для ответа
	user, err := s.repo.GetUserByID(comment.UserID)
	if err != nil {
		return nil, err
	}

	response := &dto.CommentResponse{
		ID:   comment.ID,
		Text: comment.Text,
		Date: comment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Author: dto.Author{
			Name:      user.Name,
			Surname:   user.Surname,
			AvatarURL: user.AvatarURL,
		},
	}

	return response, nil
}

// EditComment редактирует текст комментария
func (s *promoService) EditComment(userID, promoID, commentID, text string) (*dto.CommentResponse, error) {
	// Проверяем, что промокод существует
	_, err := s.repo.GetPromoByID(promoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrNotFound
		}
		return nil, err
	}

	comment, err := s.repo.GetCommentByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, dto.ErrNotFound
		}
		return nil, err
	}

	// Проверяем, что комментарий относится к данному промокоду
	if comment.PromoID != promoID {
		return nil, dto.ErrNotFound
	}

	// Проверяем, что пользователь является автором комментария
	if comment.UserID != userID {
		return nil, dto.ErrNoAccess
	}

	// Обновляем текст комментария
	comment.Text = text

	err = s.repo.UpdateComment(comment)
	if err != nil {
		return nil, err
	}

	// Получаем данные пользователя для ответа
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	response := &dto.CommentResponse{
		ID:   comment.ID,
		Text: comment.Text,
		Date: comment.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Author: dto.Author{
			Name:      user.Name,
			Surname:   user.Surname,
			AvatarURL: user.AvatarURL,
		},
	}

	return response, nil
}

// DeleteComment удаляет комментарий
func (s *promoService) DeleteComment(userID, promoID, commentID string) error {
	// Проверяем, что промокод существует
	_, err := s.repo.GetPromoByID(promoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrNotFound
		}
		return err
	}

	comment, err := s.repo.GetCommentByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ErrNotFound
		}
		return err
	}

	// Проверяем, что комментарий относится к данному промокоду
	if comment.PromoID != promoID {
		return dto.ErrNotFound
	}

	// Проверяем, что пользователь является автором комментария
	if comment.UserID != userID {
		return dto.ErrNoAccess
	}

	err = s.repo.DeleteComment(commentID)
	if err != nil {
		return err
	}

	return nil
}
