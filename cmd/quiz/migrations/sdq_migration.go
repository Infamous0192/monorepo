package migrations

import (
	"app/pkg/quiz/domain/entity"
	"log"
	"time"

	"gorm.io/gorm"
)

// SeedSDQQuiz creates the Strengths and Difficulties Questionnaire (SDQ) for children ages 4-10
func SeedSDQQuiz(db *gorm.DB) error {
	// Check if we already have the SDQ quiz
	var count int64
	if err := db.Model(&entity.Quiz{}).Where("name = ?", "SDQ Anak Usia 4-10 Tahun").Count(&count).Error; err != nil {
		return err
	}

	// If SDQ quiz already exists, skip seeding
	if count > 0 {
		log.Println("SDQ quiz already exists, skipping seed")
		return nil
	}

	log.Println("Seeding SDQ quiz data...")

	// Create SDQ quiz
	sdqQuiz := &entity.Quiz{
		Name:        "SDQ Anak Usia 4-10 Tahun",
		Description: "Strengths and Difficulties Questionnaire untuk anak usia 4-10 tahun",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := db.Create(sdqQuiz).Error; err != nil {
		return err
	}

	// Category 1: Gejala Emosional (Emotional)
	emotionalQuestions := []entity.Question{
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Sering mengeluh sakit kepala, sakit perut atau sakit-sakit lainnya.",
			Category:  "Emotional",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Banyak kekhawatiran atau sering tampak khawatir.",
			Category:  "Emotional",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Sering merasa tidak bahagia, sedih atau menangis.",
			Category:  "Emotional",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Gugup atau sulit berpisah dengan orangtua/pengasuhnya pada situasi baru, mudah kehilangan rasa percaya diri.",
			Category:  "Emotional",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Banyak yang ditakuti, mudah menjadi takut.",
			Category:  "Emotional",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := db.Create(&emotionalQuestions).Error; err != nil {
		return err
	}

	// Category 2: Masalah Perilaku (Conduct Problems)
	conductQuestions := []entity.Question{
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Sering sulit mengendalikan kemarahan.",
			Category:  "Conduct",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Umumnya bertingkah laku baik, biasanya melakukan apa yang disuruh oleh orang dewasa.",
			Category:  "Conduct",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Sering berkelahi dengan anak-anak lain atau mengintimidasi mereka.",
			Category:  "Conduct",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Sering berbohong atau berbuat curang.",
			Category:  "Conduct",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Mencuri dari rumah, sekolah atau tempat lain.",
			Category:  "Conduct",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := db.Create(&conductQuestions).Error; err != nil {
		return err
	}

	// Category 3: Hiperaktivitas (Hyperactivity)
	hyperactivityQuestions := []entity.Question{
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Gelisah, terlalu aktif, tidak dapat diam untuk waktu lama.",
			Category:  "Hyperactivity",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Terus menerus bergerak dengan resah atau menggeliat-geliat.",
			Category:  "Hyperactivity",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Mudah teralih perhatiannya, tidak dapat berkonsentrasi.",
			Category:  "Hyperactivity",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Sebelum melakukan sesuatu ia berpikir dahulu tentang akibatnya.",
			Category:  "Hyperactivity",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Memiliki perhatian yang baik terhadap apapun, mampu menyelesaikan tugas atau pekerjaan rumah sampai selesai.",
			Category:  "Hyperactivity",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := db.Create(&hyperactivityQuestions).Error; err != nil {
		return err
	}

	// Category 4: Masalah Teman Sebaya (Peer problems)
	peerQuestions := []entity.Question{
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Cenderung menyendiri, lebih suka bermain seorang diri.",
			Category:  "Peer",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Mempunyai satu atau lebih teman baik.",
			Category:  "Peer",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Pada umumnya disukai oleh anak-anak lain.",
			Category:  "Peer",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Diganggu, dipermainkan, diintimidasi atau diancam oleh anak-anak lain.",
			Category:  "Peer",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Lebih mudah berteman dengan orang dewasa dari pada dengan anak-anak lain.",
			Category:  "Peer",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := db.Create(&peerQuestions).Error; err != nil {
		return err
	}

	// Category 5: Prososial (Prosocial)
	prosocialQuestions := []entity.Question{
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Dapat memperdulikan perasaan orang lain.",
			Category:  "Prosocial",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Kalau mempunyai mainan, kesenangan atau pensil, anak bersedia berbagi dengan anak lain.",
			Category:  "Prosocial",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Suka menolong jika seseorang terluka, kecewa atau merasa sakit.",
			Category:  "Prosocial",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Bersikap baik terhadap anak-anak yang lebih muda.",
			Category:  "Prosocial",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqQuiz.ID,
			Text:      "Sering menawarkan diri untuk membantu orang lain (orang tua, guru, anak-anak lain).",
			Category:  "Prosocial",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := db.Create(&prosocialQuestions).Error; err != nil {
		return err
	}

	// Get all questions for creating answers
	var allQuestions []entity.Question
	allQuestions = append(allQuestions, emotionalQuestions...)
	allQuestions = append(allQuestions, conductQuestions...)
	allQuestions = append(allQuestions, hyperactivityQuestions...)
	allQuestions = append(allQuestions, peerQuestions...)
	allQuestions = append(allQuestions, prosocialQuestions...)

	// Create standard answer options for each question
	for _, question := range allQuestions {
		var answers []entity.Answer

		// Special inverse scoring for specific questions
		inverseScoring := false
		if (question.Category == "Conduct" && question.Text == "Umumnya bertingkah laku baik, biasanya melakukan apa yang disuruh oleh orang dewasa.") ||
			(question.Category == "Hyperactivity" && (question.Text == "Sebelum melakukan sesuatu ia berpikir dahulu tentang akibatnya." ||
				question.Text == "Memiliki perhatian yang baik terhadap apapun, mampu menyelesaikan tugas atau pekerjaan rumah sampai selesai.")) ||
			(question.Category == "Peer" && (question.Text == "Mempunyai satu atau lebih teman baik." ||
				question.Text == "Pada umumnya disukai oleh anak-anak lain.")) {
			inverseScoring = true
		}

		if inverseScoring {
			answers = []entity.Answer{
				{
					QuestionID: question.ID,
					Text:       "Tidak benar",
					Value:      2,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				},
				{
					QuestionID: question.ID,
					Text:       "Agak benar",
					Value:      1,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				},
				{
					QuestionID: question.ID,
					Text:       "Benar",
					Value:      0,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				},
			}
		} else {
			answers = []entity.Answer{
				{
					QuestionID: question.ID,
					Text:       "Tidak benar",
					Value:      0,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				},
				{
					QuestionID: question.ID,
					Text:       "Agak benar",
					Value:      1,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				},
				{
					QuestionID: question.ID,
					Text:       "Benar",
					Value:      2,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				},
			}
		}

		if err := db.Create(&answers).Error; err != nil {
			return err
		}
	}

	log.Println("SDQ quiz data seeded successfully")
	return nil
}

// SeedSDQTeenQuiz creates the Strengths and Difficulties Questionnaire (SDQ) for children ages 11-18
func SeedSDQTeenQuiz(db *gorm.DB) error {
	// Check if we already have the SDQ Teen quiz
	var count int64
	if err := db.Model(&entity.Quiz{}).Where("name = ?", "SDQ Anak Usia 11-18 Tahun").Count(&count).Error; err != nil {
		return err
	}

	// If SDQ Teen quiz already exists, skip seeding
	if count > 0 {
		log.Println("SDQ Teen quiz already exists, skipping seed")
		return nil
	}

	log.Println("Seeding SDQ Teen quiz data...")

	// Create SDQ Teen quiz
	sdqTeenQuiz := &entity.Quiz{
		Name:        "SDQ Anak Usia 11-18 Tahun",
		Description: "Strengths and Difficulties Questionnaire untuk anak usia 11-18 tahun",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := db.Create(sdqTeenQuiz).Error; err != nil {
		return err
	}

	// Category 1: Gejala Emosional (Emotional)
	emotionalQuestions := []entity.Question{
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Saya sering sakit kepala, sakit perut atau macam-macam sakit lainnya.",
			Category:  "Emotional",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Saya banyak merasa cemas atau khawatir terhadap apapun.",
			Category:  "Emotional",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Saya sering merasa tidak bahagia, sedih dan menangis.",
			Category:  "Emotional",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Saya merasa gugup dalam situasi baru. Saya sulit memusatkan perhatian pada apapun.",
			Category:  "Emotional",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Banyak yang saya takuti. Saya mudah menjadi takut.",
			Category:  "Emotional",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := db.Create(&emotionalQuestions).Error; err != nil {
		return err
	}

	// Category 2: Masalah Perilaku (Conduct Problems)
	conductQuestions := []entity.Question{
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Saya menjadi sangat marah dan sering tidak bisa mengendalikan kemarahan saya.",
			Category:  "Conduct",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Saya biasanya melakukan apa yang diperintahkan oleh orang lain.",
			Category:  "Conduct",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Saya sering bertengkar dengan orang lain. Saya dapat memaksa orang lain untuk melakukan apa yang saya inginkan.",
			Category:  "Conduct",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Saya sering dituduh berbohong atau berbuat curang.",
			Category:  "Conduct",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Saya mengambil barang yang bukan milik saya dari rumah, sekolah, atau darimana saja.",
			Category:  "Conduct",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := db.Create(&conductQuestions).Error; err != nil {
		return err
	}

	// Category 3: Hiperaktivitas (Hyperactivity)
	hyperactivityQuestions := []entity.Question{
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Saya gelisah, saya tidak dapat diam untuk waktu lama.",
			Category:  "Hyperactivity",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Bila sedang gelisah atau cemas badan saya sering bergerak-gerak tanpa saya sadari.",
			Category:  "Hyperactivity",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Perhatian saya mudah teralihkan. Saya sulit memusatkan perhatian pada apapun.",
			Category:  "Hyperactivity",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Sebelum melakukan sesuatu saya berpikir dahulu tentang akibatnya.",
			Category:  "Hyperactivity",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Saya menyelesaikan pekerjaan yang sedang saya lakukan. Saya mempunyai perhatian yang baik terhadap apapun.",
			Category:  "Hyperactivity",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := db.Create(&hyperactivityQuestions).Error; err != nil {
		return err
	}

	// Category 4: Masalah Teman Sebaya (Peer problems)
	peerQuestions := []entity.Question{
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Saya lebih suka sendirian dari pada bersama dengan orang-orang yang seumuran saya.",
			Category:  "Peer",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Saya mempunyai satu teman baik atau lebih.",
			Category:  "Peer",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Orang lain seumur saya pada umumnya menyukai saya.",
			Category:  "Peer",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Saya sering diganggu atau dipermainkan oleh anak-anak atau remaja lainnya.",
			Category:  "Peer",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Saya lebih mudah berteman dengan orang dewasa daripada dengan orang-orang seumuran saya.",
			Category:  "Peer",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := db.Create(&peerQuestions).Error; err != nil {
		return err
	}

	// Category 5: Prososial (Prosocial)
	prosocialQuestions := []entity.Question{
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Saya berusaha bersikap baik kepada orang lain. Saya peduli dengan perasaan mereka.",
			Category:  "Prosocial",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Kalau saya memiliki mainan, CD, atau makanan saya biasanya berbagi dengan orang lain.",
			Category:  "Prosocial",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Saya selalu siap menolong jika ada orang terluka, kecewa atau merasa sakit.",
			Category:  "Prosocial",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Saya bersikap baik pada anak-anak yang lebih muda dari saya.",
			Category:  "Prosocial",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    sdqTeenQuiz.ID,
			Text:      "Saya sering menawarkan diri untuk membantu orang lain, orang tua, guru atau anak-anak.",
			Category:  "Prosocial",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := db.Create(&prosocialQuestions).Error; err != nil {
		return err
	}

	// Get all questions for creating answers
	var allQuestions []entity.Question
	allQuestions = append(allQuestions, emotionalQuestions...)
	allQuestions = append(allQuestions, conductQuestions...)
	allQuestions = append(allQuestions, hyperactivityQuestions...)
	allQuestions = append(allQuestions, peerQuestions...)
	allQuestions = append(allQuestions, prosocialQuestions...)

	// Create standard answer options for each question
	for _, question := range allQuestions {
		var answers []entity.Answer

		// Special inverse scoring for specific questions
		inverseScoring := false
		if (question.Category == "Conduct" && question.Text == "Saya biasanya melakukan apa yang diperintahkan oleh orang lain.") ||
			(question.Category == "Hyperactivity" && (question.Text == "Sebelum melakukan sesuatu saya berpikir dahulu tentang akibatnya." ||
				question.Text == "Saya menyelesaikan pekerjaan yang sedang saya lakukan. Saya mempunyai perhatian yang baik terhadap apapun.")) ||
			(question.Category == "Peer" && (question.Text == "Saya mempunyai satu teman baik atau lebih." ||
				question.Text == "Orang lain seumur saya pada umumnya menyukai saya.")) {
			inverseScoring = true
		}

		if inverseScoring {
			answers = []entity.Answer{
				{
					QuestionID: question.ID,
					Text:       "Tidak benar",
					Value:      2,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				},
				{
					QuestionID: question.ID,
					Text:       "Agak benar",
					Value:      1,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				},
				{
					QuestionID: question.ID,
					Text:       "Benar",
					Value:      0,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				},
			}
		} else {
			answers = []entity.Answer{
				{
					QuestionID: question.ID,
					Text:       "Tidak benar",
					Value:      0,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				},
				{
					QuestionID: question.ID,
					Text:       "Agak benar",
					Value:      1,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				},
				{
					QuestionID: question.ID,
					Text:       "Benar",
					Value:      2,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				},
			}
		}

		if err := db.Create(&answers).Error; err != nil {
			return err
		}
	}

	log.Println("SDQ Teen quiz data seeded successfully")
	return nil
}
