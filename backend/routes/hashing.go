package routes
import(
	"golang.org/x/crypto/bcrypt"
)
func HashPassword(password string) (string, error){
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err

}
func CheckPasswordHash(password, hash string) bool{
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return  err == nil
}