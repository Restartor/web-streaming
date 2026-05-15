package repository

import (
	"backend/internal/domain"

	"gorm.io/gorm"
)

// RefreshTokenRepository mengelola penyimpanan refresh token di database.
// Refresh token adalah kredensial jangka panjang yang terikat pada User (melalui UserID) yang memungkinkan
// klien untuk mendapatkan token akses baru tanpa memasukkan kembali kredensial.
// Setiap record RefreshToken menyimpan string token unik yang terikat pada user tertentu,
// memungkinkan manajemen sesi stateful dan rotasi token yang aman.
type RefreshTokenRepository struct {
	db *gorm.DB
}

// Create menyimpan refresh token baru di database untuk user tertentu.
// Dipanggil saat login untuk mengasosiasikan token jangka panjang dengan sesi user.
func (r *RefreshTokenRepository) Create(token *domain.RefreshToken) error {
	return r.db.Create(token).Error
}

// FindByToken mengambil refresh token berdasarkan nilai stringnya.
// Digunakan ketika klien mengirimkan refresh token untuk memverifikasi bahwa token tersebut ada dan milik user yang terautentikasi.
func (r *RefreshTokenRepository) FindByToken(token string) (*domain.RefreshToken, error) {
	var refreshToken domain.RefreshToken
	err := r.db.Where("token = ?", token).First(&refreshToken).Error
	return &refreshToken, err
}

// DeleteByUserID menghapus semua refresh token untuk user tertentu (UserID).
// Dipanggil saat logout atau reset password untuk membatalkan semua sesi aktif user tersebut.
func (r *RefreshTokenRepository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&domain.RefreshToken{}).Error
}

func (r *RefreshTokenRepository) DeleteByToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&domain.RefreshToken{}).Error
}

// NewRefreshTokenRepository membuat instance baru dari RefreshTokenRepository.
// Menghubungkan interface domain.RefreshTokenRepository dengan implementasi database konkretnya.
// Ini memungkinkan UserService untuk mengelola refresh token sambil mengabstraksi layer database.
func NewRefreshTokenRepository(db *gorm.DB) domain.RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}
