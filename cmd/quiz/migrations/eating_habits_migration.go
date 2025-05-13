package migrations

import (
	"app/pkg/quiz/domain/entity"
	"log"
	"time"

	"gorm.io/gorm"
)

// SeedEatingHabitsQuiz creates the Eating Habits Questionnaire for children
func SeedEatingHabitsQuiz(db *gorm.DB) error {
	// Check if we already have the Eating Habits quiz
	var count int64
	if err := db.Model(&entity.Quiz{}).Where("name = ?", "Kuesioner Kebiasaan Makan Anak").Count(&count).Error; err != nil {
		return err
	}

	// If Eating Habits quiz already exists, skip seeding
	if count > 0 {
		log.Println("Eating Habits quiz already exists, skipping seed")
		return nil
	}

	log.Println("Seeding Eating Habits quiz data...")

	// Create Eating Habits quiz
	eatingHabitsQuiz := &entity.Quiz{
		Name:        "Kuesioner Kebiasaan Makan Anak",
		Description: "Kuesioner untuk menilai kebiasaan makan dan pola gizi anak",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := db.Create(eatingHabitsQuiz).Error; err != nil {
		return err
	}

	// Add all sections and questions
	if err := addEatingHabitsQuestions(db, eatingHabitsQuiz.ID); err != nil {
		return err
	}

	if err := addDrinkingHabitsQuestions(db, eatingHabitsQuiz.ID); err != nil {
		return err
	}

	if err := addPhysicalActivityQuestions(db, eatingHabitsQuiz.ID); err != nil {
		return err
	}

	if err := addNutritionAssessmentQuestions(db, eatingHabitsQuiz.ID); err != nil {
		return err
	}

	log.Println("Eating Habits quiz data seeded successfully")
	return nil
}

// addEatingHabitsQuestions adds the eating habits section questions
func addEatingHabitsQuestions(db *gorm.DB, quizID uint) error {
	eatingHabitsQuestions := []entity.Question{
		{
			QuizID:    quizID,
			Text:      "Berapa kali kamu makan dalam sehari?",
			Category:  "Kebiasaan Makan",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Apakah kamu sarapan atau makan pagi sebelum berangkat ke sekolah?",
			Category:  "Kebiasaan Makan",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Apa makanan favorit kamu?",
			Category:  "Kebiasaan Makan",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Apakah kamu suka makan sayur?",
			Category:  "Kebiasaan Makan",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Apakah kamu suka mencoba makanan baru?",
			Category:  "Kebiasaan Makan",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Seberapa sering kamu makan camilan dalam sehari?",
			Category:  "Kebiasaan Makan",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Camilan apa yang kamu sukai?",
			Category:  "Kebiasaan Makan",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := db.Create(&eatingHabitsQuestions).Error; err != nil {
		return err
	}

	// Create answers for question 1: Frequency of meals
	mealFrequencyAnswers := []entity.Answer{
		{
			QuestionID: eatingHabitsQuestions[0].ID,
			Text:       "1 kali",
			Value:      0, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: eatingHabitsQuestions[0].ID,
			Text:       "2 kali",
			Value:      1, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: eatingHabitsQuestions[0].ID,
			Text:       "3 kali",
			Value:      2, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: eatingHabitsQuestions[0].ID,
			Text:       "Lebih dari 3 kali",
			Value:      2, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&mealFrequencyAnswers).Error; err != nil {
		return err
	}

	// Create answers for question 2: Breakfast
	breakfastAnswers := []entity.Answer{
		{
			QuestionID: eatingHabitsQuestions[1].ID,
			Text:       "Ya, selalu",
			Value:      2, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: eatingHabitsQuestions[1].ID,
			Text:       "Kadang-kadang",
			Value:      1, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: eatingHabitsQuestions[1].ID,
			Text:       "Tidak pernah",
			Value:      0, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&breakfastAnswers).Error; err != nil {
		return err
	}

	// Create answers for question 3: Favorite food
	favoriteAnswers := []entity.Answer{
		{
			QuestionID: eatingHabitsQuestions[2].ID,
			Text:       "Sayur-sayuran",
			Value:      2, // Healthier choice
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: eatingHabitsQuestions[2].ID,
			Text:       "Buah-buahan",
			Value:      2, // Healthier choice
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: eatingHabitsQuestions[2].ID,
			Text:       "Nasi atau roti",
			Value:      1, // Staple food
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: eatingHabitsQuestions[2].ID,
			Text:       "Daging atau ikan",
			Value:      1, // Protein source
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: eatingHabitsQuestions[2].ID,
			Text:       "Camilan (permen, keripik, dll)",
			Value:      0, // Less healthy choice
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&favoriteAnswers).Error; err != nil {
		return err
	}

	// Create answers for question 4: Vegetable consumption
	vegetableAnswers := []entity.Answer{
		{
			QuestionID: eatingHabitsQuestions[3].ID,
			Text:       "Selalu",
			Value:      2, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: eatingHabitsQuestions[3].ID,
			Text:       "Kadang-kadang",
			Value:      1, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: eatingHabitsQuestions[3].ID,
			Text:       "Jarang sekali",
			Value:      0, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&vegetableAnswers).Error; err != nil {
		return err
	}

	// Create answers for question 5: Trying new foods
	newFoodAnswers := []entity.Answer{
		{
			QuestionID: eatingHabitsQuestions[4].ID,
			Text:       "Selalu",
			Value:      2, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: eatingHabitsQuestions[4].ID,
			Text:       "Kadang-kadang",
			Value:      1, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: eatingHabitsQuestions[4].ID,
			Text:       "Jarang sekali",
			Value:      0, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&newFoodAnswers).Error; err != nil {
		return err
	}

	// Create answers for question 6: Snack frequency
	snackFrequencyAnswers := []entity.Answer{
		{
			QuestionID: eatingHabitsQuestions[5].ID,
			Text:       "Tidak pernah",
			Value:      2, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: eatingHabitsQuestions[5].ID,
			Text:       "Sekali",
			Value:      2, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: eatingHabitsQuestions[5].ID,
			Text:       "Dua kali",
			Value:      1, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: eatingHabitsQuestions[5].ID,
			Text:       "Lebih dari dua kali",
			Value:      0, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&snackFrequencyAnswers).Error; err != nil {
		return err
	}

	// Create answers for question 7: Preferred snacks
	snackTypeAnswers := []entity.Answer{
		{
			QuestionID: eatingHabitsQuestions[6].ID,
			Text:       "Buah-buahan",
			Value:      1, // +1 for healthy snack
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: eatingHabitsQuestions[6].ID,
			Text:       "Keripik atau snack kemasan",
			Value:      -1, // -1 for unhealthy snack
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: eatingHabitsQuestions[6].ID,
			Text:       "Permen atau cokelat",
			Value:      -1, // -1 for unhealthy snack
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: eatingHabitsQuestions[6].ID,
			Text:       "Kue atau roti manis",
			Value:      -1, // -1 for unhealthy snack
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: eatingHabitsQuestions[6].ID,
			Text:       "Kacang-kacangan, yogurt",
			Value:      1, // +1 for healthy snack
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&snackTypeAnswers).Error; err != nil {
		return err
	}

	return nil
}

// addDrinkingHabitsQuestions adds the drinking habits section questions
func addDrinkingHabitsQuestions(db *gorm.DB, quizID uint) error {
	drinkingQuestions := []entity.Question{
		{
			QuizID:    quizID,
			Text:      "Apa minuman yang paling sering kamu konsumsi?",
			Category:  "Minuman",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Berapa banyak air putih yang kamu minum dalam sehari?",
			Category:  "Minuman",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := db.Create(&drinkingQuestions).Error; err != nil {
		return err
	}

	// Create answers for question 1: Preferred drink
	drinkTypeAnswers := []entity.Answer{
		{
			QuestionID: drinkingQuestions[0].ID,
			Text:       "Air putih",
			Value:      2, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: drinkingQuestions[0].ID,
			Text:       "Susu",
			Value:      2, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: drinkingQuestions[0].ID,
			Text:       "Jus buah",
			Value:      1, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: drinkingQuestions[0].ID,
			Text:       "Minuman manis (soft drink)",
			Value:      0, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&drinkTypeAnswers).Error; err != nil {
		return err
	}

	// Create answers for question 2: Water consumption
	waterConsumptionAnswers := []entity.Answer{
		{
			QuestionID: drinkingQuestions[1].ID,
			Text:       "Sedikit",
			Value:      0, // Less healthy
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: drinkingQuestions[1].ID,
			Text:       "2-3 gelas",
			Value:      1, // Moderate
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: drinkingQuestions[1].ID,
			Text:       "Banyak lebih dari 3 gelas",
			Value:      2, // Healthier
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&waterConsumptionAnswers).Error; err != nil {
		return err
	}

	return nil
}

// addPhysicalActivityQuestions adds the physical activity section questions
func addPhysicalActivityQuestions(db *gorm.DB, quizID uint) error {
	activityQuestion := entity.Question{
		QuizID:    quizID,
		Text:      "Seberapa sering kamu bermain atau berolahraga dalam satu minggu?",
		Category:  "Aktivitas Fisik",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&activityQuestion).Error; err != nil {
		return err
	}

	// Create answers for physical activity frequency
	activityAnswers := []entity.Answer{
		{
			QuestionID: activityQuestion.ID,
			Text:       "Setiap hari",
			Value:      2, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: activityQuestion.ID,
			Text:       "2-4 kali seminggu",
			Value:      1, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: activityQuestion.ID,
			Text:       "Kadang-kadang",
			Value:      0, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: activityQuestion.ID,
			Text:       "Jarang sekali",
			Value:      0, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&activityAnswers).Error; err != nil {
		return err
	}

	return nil
}

// addNutritionAssessmentQuestions adds the nutrition assessment section questions
func addNutritionAssessmentQuestions(db *gorm.DB, quizID uint) error {
	nutritionQuestions := []entity.Question{
		{
			QuizID:    quizID,
			Text:      "Apakah kamu merasa kenyang setelah makan?",
			Category:  "Penilaian Gizi",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Siapa yang biasanya menyiapkan makanan untukmu?",
			Category:  "Penilaian Gizi",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := db.Create(&nutritionQuestions).Error; err != nil {
		return err
	}

	// Create answers for question 1: Fullness after eating
	fullnessAnswers := []entity.Answer{
		{
			QuestionID: nutritionQuestions[0].ID,
			Text:       "Selalu kenyang",
			Value:      2, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: nutritionQuestions[0].ID,
			Text:       "Kadang-kadang kenyang",
			Value:      1, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: nutritionQuestions[0].ID,
			Text:       "Tidak merasa kenyang",
			Value:      0, // Based on scoring system
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&fullnessAnswers).Error; err != nil {
		return err
	}

	// Create answers for question 2: Food preparation
	prepAnswers := []entity.Answer{
		{
			QuestionID: nutritionQuestions[1].ID,
			Text:       "Ibu atau ayah",
			Value:      2, // Parents typically provide more balanced meals
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: nutritionQuestions[1].ID,
			Text:       "Kakek atau nenek",
			Value:      2, // Grandparents typically provide more balanced meals
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: nutritionQuestions[1].ID,
			Text:       "Kakak atau adik",
			Value:      1, // Siblings may provide less balanced meals
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: nutritionQuestions[1].ID,
			Text:       "Saya sendiri",
			Value:      0, // Children may choose less balanced meals
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&prepAnswers).Error; err != nil {
		return err
	}

	return nil
}
