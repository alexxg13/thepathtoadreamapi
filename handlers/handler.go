package handlers

import (
	"10042025/AI"
	"10042025/pkg/postgreSQL"
	"github.com/coalaura/mistral"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
)

const (
	usersURL    = "/users/"
	loginURL    = "/login/"
	registerURL = "/registration/"
	userURL     = "/user/"
	chatURL     = "/chat/"
)

type handler struct {
	db       *gorm.DB
	clientAI *mistral.MistralClient
}

// NewHandler создает новый обработчик с подключением к базе данных и AI-клиенту.
func NewHandler(DB *gorm.DB, client *mistral.MistralClient) Handler {
	return &handler{
		db:       DB,
		clientAI: client,
	}
}

// Register регистрирует все маршруты для авторизации, пользователей и чатов.
func (h *handler) Register(router *gin.Engine) {
	// Регистрация и авторизация
	router.POST(registerURL, h.registration)
	router.POST(loginURL, h.login)

	// Работа с пользователями
	router.GET(userURL, h.getUser)
	router.GET(usersURL, h.getUsers)

	// Работа с диалогами AI
	//router.GET(historyURL, h.getHistory)

	// Работа с Чатами и Роадмап
	router.POST(chatURL, h.sendPrompt)
}

// login - авторизация пользователя
// @Summary      Авторизация пользователя
// @Description  Проверяет учетные данные пользователя и возвращает статус входа
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body object{username=string,password=string} true "Данные для входа"
// @Success      200 {object} map[string]string "User logged successfully"
// @Failure      400 {object} map[string]string "Invalid request data"
// @Failure      401 {object} map[string]string "Invalid credentials"
// @Router       /login/ [post]
func (h *handler) login(c *gin.Context) {
	var user postgreSQL.User
	var jsonMessage struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}

	// Проверка входных данных
	if err := c.ShouldBindJSON(&jsonMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Поиск пользователя в базе данных
	if err := h.db.Where("username = ?", jsonMessage.UserName).First(&user).Error; err != nil {
		log.Printf("Error login: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"username": jsonMessage.UserName,
			"error":    "Invalid credentials",
		})
		return
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(jsonMessage.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Успешная авторизация
	c.JSON(http.StatusOK, gin.H{"message": "User logged successfully"})
}

// registration - регистрация пользователя
// @Summary      Регистрация пользователя
// @Description  Регистрирует нового пользователя в системе
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body object{username=string,password=string} true "Данные нового пользователя"
// @Success      200 {object} map[string]string "User registered successfully"
// @Failure      400 {object} map[string]string "Invalid request data"
// @Failure      409 {object} map[string]string "User already exists"
// @Router       /registration/ [post]
func (h *handler) registration(c *gin.Context) {
	var user postgreSQL.User

	// Получение данных пользователя из запроса
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка на существующего пользователя
	searchUsername := h.db.Where("username = ?", user.Username).First(&user).Error
	if searchUsername == nil {
		log.Printf("Error existing user")
		c.JSON(http.StatusConflict, gin.H{"message": "User already exists"})
		return
	} else if searchUsername != gorm.ErrRecordNotFound {
		log.Printf("Error search user: %s", searchUsername)
		c.JSON(http.StatusConflict, gin.H{"message": "Error search user"})
		return
	}

	// Хэширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}
	user.PasswordHash = string(hashedPassword)

	// Сохранение пользователя в БД
	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// getUser - получение информации о пользователе по имени
// @Summary      Получение информации о пользователе
// @Description  Возвращает информацию о пользователе по имени
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        username query string true "Имя пользователя"
// @Success      200 {object} map[string]interface{} "User information"
// @Failure      500 {object} map[string]string "Error search User"
// @Router       /user/ [get]
func (h *handler) getUser(c *gin.Context) {
	var user postgreSQL.User
	username := c.Query("username")

	// Поиск пользователя в базе данных
	if err := h.db.Where("username = ?", username).Preload("Dialogs").First(&user).Error; err != nil {
		log.Printf("Error GetUser: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error search User"})
		return
	}

	// Возвращаем информацию о пользователе
	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		"Dialogs":  user.Dialogs,
	})
}

// getUsers - получение списка всех пользователей и их диалогов
// @Summary      Получение списка пользователей
// @Description  Возвращает список всех пользователей и их диалоги
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200 {array} map[string]interface{} "List of users"
// @Failure      400 {object} map[string]string "Error find Users"
// @Router       /users/ [get]
func (h *handler) getUsers(c *gin.Context) {
	var users []postgreSQL.User

	// Получаем список пользователей из базы данных
	if err := h.db.Preload("Dialogs").Find(&users).Error; err != nil {
		log.Printf("Error find Users: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error find Users"})
		return
	}

	// Возвращаем список пользователей
	c.JSON(http.StatusOK, gin.H{"message": users})
}

// sendPrompt - отправка запроса в AI и сохранение диалога
// @Summary      Отправка запроса в AI
// @Description  Отправляет данные в AI и сохраняет диалог в БД
// @Tags         ai
// @Accept       json
// @Produce      json
// @Param        request body object{userID=uint,goal=string,message=[]string} true "Запрос AI, message содержит 8 значений!"
// @Success      200 {object} map[string]string "AI response"
// @Failure      400 {object} map[string]string "Invalid request data"
// @Failure      500 {object} map[string]string "Error response Mistral AI"
// @Router       /chat/ [post]
func (h *handler) sendPrompt(c *gin.Context) {
	var chatData struct {
		UserID  uint     `json:"userID"`
		Goal    string   `json:"goal"`
		Message []string `json:"message"`
	}

	// Проверка данных запроса
	if err := c.ShouldBindJSON(&chatData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(chatData.Message) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error len message < 8"})
		return
	}

	// Отправка запроса в AI
	response, err := AI.SendPrompt(h.clientAI, chatData.Goal, chatData.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error response Mistral AI",
		})
		return
	}

	response = strings.Replace(response, "#", "", -1)
	response = strings.Replace(response, "*", "", -1)

	// Сохранение диалога в БД
	dialog := postgreSQL.Dialog{
		UserID:   chatData.UserID,
		Message:  AI.GenerationPromt(chatData.Goal, chatData.Message),
		Response: response,
	}
	if err := h.db.Create(&dialog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save dialog"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": response})
}
