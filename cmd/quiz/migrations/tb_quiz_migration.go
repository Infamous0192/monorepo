package migrations

import (
	"app/pkg/quiz/domain/entity"
	"log"
	"time"

	"gorm.io/gorm"
)

// SeedTBQuiz creates the Tuberculosis Knowledge, Attitude, and Practice (KAP) quiz
func SeedTBQuiz(db *gorm.DB) error {
	// Check if we already have the TB quiz
	var count int64
	if err := db.Model(&entity.Quiz{}).Where("name = ?", "Kuesioner Tuberkulosis (TB)").Count(&count).Error; err != nil {
		return err
	}

	// If TB quiz already exists, skip seeding
	if count > 0 {
		log.Println("TB quiz already exists, skipping seed")
		return nil
	}

	log.Println("Seeding TB quiz data...")

	// Create TB quiz
	tbQuiz := &entity.Quiz{
		Name:        "Kuesioner Tuberkulosis (TB)",
		Description: "Kuesioner untuk menilai pengetahuan, sikap, dan praktik terkait Tuberkulosis (TB)",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := db.Create(tbQuiz).Error; err != nil {
		return err
	}

	// Add all sections and questions
	if err := addSociodemographicQuestions(db, tbQuiz.ID); err != nil {
		return err
	}

	if err := addKnowledgeQuestions(db, tbQuiz.ID); err != nil {
		return err
	}

	if err := addAttitudeQuestions(db, tbQuiz.ID); err != nil {
		return err
	}

	if err := addHousingConditionQuestions(db, tbQuiz.ID); err != nil {
		return err
	}

	if err := addSmokingBehaviorQuestions(db, tbQuiz.ID); err != nil {
		return err
	}

	log.Println("TB quiz data seeded successfully")
	return nil
}

// addSociodemographicQuestions adds the sociodemographic section questions
func addSociodemographicQuestions(db *gorm.DB, quizID uint) error {
	sociodemographicQuestions := []entity.Question{
		{
			QuizID:    quizID,
			Text:      "Jenis Kelamin",
			Category:  "Sosiodemografis",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Pendidikan",
			Category:  "Sosiodemografis",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Pekerjaan",
			Category:  "Sosiodemografis",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Penghasilan",
			Category:  "Sosiodemografis",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Apakah memiliki anggota keluarga yang punya kebiasaan merokok?",
			Category:  "Sosiodemografis",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := db.Create(&sociodemographicQuestions).Error; err != nil {
		return err
	}

	// Create answers for each question
	genderAnswers := []entity.Answer{
		{
			QuestionID: sociodemographicQuestions[0].ID,
			Text:       "Laki-laki",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: sociodemographicQuestions[0].ID,
			Text:       "Perempuan",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&genderAnswers).Error; err != nil {
		return err
	}

	educationAnswers := []entity.Answer{
		{
			QuestionID: sociodemographicQuestions[1].ID,
			Text:       "Tidak Sekolah",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: sociodemographicQuestions[1].ID,
			Text:       "SD",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: sociodemographicQuestions[1].ID,
			Text:       "SMP",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: sociodemographicQuestions[1].ID,
			Text:       "SMA",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: sociodemographicQuestions[1].ID,
			Text:       "Perguruan Tinggi",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&educationAnswers).Error; err != nil {
		return err
	}

	employmentAnswers := []entity.Answer{
		{
			QuestionID: sociodemographicQuestions[2].ID,
			Text:       "Tidak Bekerja",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: sociodemographicQuestions[2].ID,
			Text:       "Bekerja",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&employmentAnswers).Error; err != nil {
		return err
	}

	incomeAnswers := []entity.Answer{
		{
			QuestionID: sociodemographicQuestions[3].ID,
			Text:       "Kurang dari Rp3.500.000",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: sociodemographicQuestions[3].ID,
			Text:       "Lebih dari Rp3.500.000",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&incomeAnswers).Error; err != nil {
		return err
	}

	smokingFamilyAnswers := []entity.Answer{
		{
			QuestionID: sociodemographicQuestions[4].ID,
			Text:       "Ya, Merokok lebih dari 6 bulan",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: sociodemographicQuestions[4].ID,
			Text:       "Tidak merokok / merokok kurang dari 6 bulan",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&smokingFamilyAnswers).Error; err != nil {
		return err
	}

	return nil
}

// addKnowledgeQuestions adds the knowledge assessment section questions
func addKnowledgeQuestions(db *gorm.DB, quizID uint) error {
	knowledgeQuestions := []entity.Question{
		{
			QuizID:    quizID,
			Text:      "Apa yang dimaksud dengan penyakit TB Paru?",
			Category:  "Pengetahuan",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Apakah yang menjadi penyebab penyakit TB Paru?",
			Category:  "Pengetahuan",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Apa saja gejala atau tanda-tanda pada penderita TB Paru?",
			Category:  "Pengetahuan",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Bagaimana cara TB Paru menular kepada seseorang?",
			Category:  "Pengetahuan",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Pengobatan TB Paru membutuhkan waktu berapa lama untuk sembuh?",
			Category:  "Pengetahuan",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Obat TB paru yang tidak diminum secara teratur hingga habis akan mengakibatkan?",
			Category:  "Pengetahuan",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Bagaimana cara terbaik untuk menghindari penularan penyakit TB Paru terhadap orang lain?",
			Category:  "Pengetahuan",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Bagaimana kondisi lingkungan atau rumah yang memicu penyakit TB Paru?",
			Category:  "Pengetahuan",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Kapankah seorang penderita TB Paru dinyatakan sembuh?",
			Category:  "Pengetahuan",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Kebiasaan seperti apa yang dapat memperburuk penyakit TB Paru?",
			Category:  "Pengetahuan",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := db.Create(&knowledgeQuestions).Error; err != nil {
		return err
	}

	// Create answers for each knowledge question
	question1Answers := []entity.Answer{
		{
			QuestionID: knowledgeQuestions[0].ID,
			Text:       "Batuk-batuk selama 3 minggu dan nyeri dada",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: knowledgeQuestions[0].ID,
			Text:       "Batuk dengan gatal ditenggorokan",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: knowledgeQuestions[0].ID,
			Text:       "Batuk-batuk akibat merokok",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&question1Answers).Error; err != nil {
		return err
	}

	question2Answers := []entity.Answer{
		{
			QuestionID: knowledgeQuestions[1].ID,
			Text:       "Bakteri",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: knowledgeQuestions[1].ID,
			Text:       "Virus",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: knowledgeQuestions[1].ID,
			Text:       "Genetik/keturunan",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&question2Answers).Error; err != nil {
		return err
	}

	question3Answers := []entity.Answer{
		{
			QuestionID: knowledgeQuestions[2].ID,
			Text:       "Batuk rejan",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: knowledgeQuestions[2].ID,
			Text:       "Batuk berdahak lebih dari 3 minggu",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: knowledgeQuestions[2].ID,
			Text:       "Batuk tidak berdahak lebih dari 3 minggu",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&question3Answers).Error; err != nil {
		return err
	}

	question4Answers := []entity.Answer{
		{
			QuestionID: knowledgeQuestions[3].ID,
			Text:       "Melalui kontak langsung (misal jabat tangan dan lain-lain)",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: knowledgeQuestions[3].ID,
			Text:       "Melalui makanan dan minuman",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: knowledgeQuestions[3].ID,
			Text:       "Melalui percikan dahak/ludah",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&question4Answers).Error; err != nil {
		return err
	}

	question5Answers := []entity.Answer{
		{
			QuestionID: knowledgeQuestions[4].ID,
			Text:       "Sepanjang hidupnya",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: knowledgeQuestions[4].ID,
			Text:       "6 bulan atau lebih setelah berobat teratur hingga tuntas",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: knowledgeQuestions[4].ID,
			Text:       "1 bulan setelah berobat tidak teratur",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&question5Answers).Error; err != nil {
		return err
	}

	question6Answers := []entity.Answer{
		{
			QuestionID: knowledgeQuestions[5].ID,
			Text:       "Tidak ada akibat",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: knowledgeQuestions[5].ID,
			Text:       "Penyakit TB sembuh secara alami",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: knowledgeQuestions[5].ID,
			Text:       "Kuman akan kebal terhadap obat sehingga tidak sembuh",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&question6Answers).Error; err != nil {
		return err
	}

	question7Answers := []entity.Answer{
		{
			QuestionID: knowledgeQuestions[6].ID,
			Text:       "Menutup mulut dan hidung saat batuk atau bersin dan tidak meludah di sembarang tempat",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: knowledgeQuestions[6].ID,
			Text:       "Batuk atau bersin",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: knowledgeQuestions[6].ID,
			Text:       "Meludah di sembarang tempat",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&question7Answers).Error; err != nil {
		return err
	}

	question8Answers := []entity.Answer{
		{
			QuestionID: knowledgeQuestions[7].ID,
			Text:       "Banyak sampah dan lembab",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: knowledgeQuestions[7].ID,
			Text:       "Pencahayaan yang baik",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: knowledgeQuestions[7].ID,
			Text:       "Memiliki banyak ventilasi",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&question8Answers).Error; err != nil {
		return err
	}

	question9Answers := []entity.Answer{
		{
			QuestionID: knowledgeQuestions[8].ID,
			Text:       "Saat batuk sudah hilang",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: knowledgeQuestions[8].ID,
			Text:       "Sampai dinyatakan sembuh oleh dokter",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: knowledgeQuestions[8].ID,
			Text:       "Saat obatnya sudah habis",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&question9Answers).Error; err != nil {
		return err
	}

	question10Answers := []entity.Answer{
		{
			QuestionID: knowledgeQuestions[9].ID,
			Text:       "Beraktivitas fisik secara teratur",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: knowledgeQuestions[9].ID,
			Text:       "Makan makanan yang bergizi",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: knowledgeQuestions[9].ID,
			Text:       "Merokok dan minum minuman beralkohol",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&question10Answers).Error; err != nil {
		return err
	}

	return nil
}

// addAttitudeQuestions adds the attitude assessment section questions
func addAttitudeQuestions(db *gorm.DB, quizID uint) error {
	attitudeQuestions := []entity.Question{
		{
			QuizID:    quizID,
			Text:      "Penyakit TB paru merupakan penyakit yang sangat menular, namun dapat disembuhkan.",
			Category:  "Sikap",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Setiap orang yang mengalami batuk berdahak selama 2 minggu atau lebih sebaiknya melakukan pemeriksaan dahak ke pelayanan kesehatan.",
			Category:  "Sikap",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Penularan TB paru dapat dicegah apabila penderita menggunakan masker dan tidak berbicara terlalu dekat dengan lawan bicaranya.",
			Category:  "Sikap",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Penderita TB paru harus menutup mulut saat batuk atau bersin untuk mencegah penyebaran kuman kepada orang lain.",
			Category:  "Sikap",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Penderita TB paru harus melakukan pengobatan secara rutin selama 6 bulan dan sampai sembuh.",
			Category:  "Sikap",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Adanya lingkungan yang bersih dan tidak padat penghuni dapat mencegah penularan TB paru.",
			Category:  "Sikap",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Membuang ludah sembarangan dapat meningkatkan risiko penularan penyakit TB paru.",
			Category:  "Sikap",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Sinar matahari dapat membunuh kuman yang ada di dalam rumah.",
			Category:  "Sikap",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Luas ventilasi yang memenuhi syarat yaitu minimal 10% dari luas lantai.",
			Category:  "Sikap",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := db.Create(&attitudeQuestions).Error; err != nil {
		return err
	}

	// Create common answers for all attitude questions
	for _, question := range attitudeQuestions {
		attitudeAnswers := []entity.Answer{
			{
				QuestionID: question.ID,
				Text:       "Setuju",
				Value:      1,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			{
				QuestionID: question.ID,
				Text:       "Tidak Setuju",
				Value:      0,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		}

		if err := db.Create(&attitudeAnswers).Error; err != nil {
			return err
		}
	}

	return nil
}

// addHousingConditionQuestions adds the housing condition section questions
func addHousingConditionQuestions(db *gorm.DB, quizID uint) error {
	housingQuestion := entity.Question{
		QuizID:    quizID,
		Text:      "Bagaimana pencahayaan atau sinar matahari yang masuk rumah, Apakah memerlukan alat penerangan seperti lampu untuk membaca buku atau koran pada siang hari di dalam rumah?",
		Category:  "Kondisi Tempat Tinggal",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&housingQuestion).Error; err != nil {
		return err
	}

	housingAnswers := []entity.Answer{
		{
			QuestionID: housingQuestion.ID,
			Text:       "Ya, memerlukan alat penerangan lampu",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: housingQuestion.ID,
			Text:       "Tidak, karena dapat membaca buku dengan jelas",
			Value:      2, // Note: value is 2 here based on the provided questionnaire
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&housingAnswers).Error; err != nil {
		return err
	}

	return nil
}

// addSmokingBehaviorQuestions adds the smoking behavior section questions
func addSmokingBehaviorQuestions(db *gorm.DB, quizID uint) error {
	smokingQuestions := []entity.Question{
		{
			QuizID:    quizID,
			Text:      "Apakah bapak/ibu mempunyai kebiasaan merokok/mudah terpapar asap rokok?",
			Category:  "Perilaku Merokok",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Berapa umur bapak/ibu ketika pertama kali merokok?",
			Category:  "Perilaku Merokok",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Berapa rata-rata batang rokok yang dihisap dalam sehari?",
			Category:  "Perilaku Merokok",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Berapa lama durasi merokok bapak/ibu dimulai dari usia awal merokok sampai pada saat penelitian dilakukan atau berhenti merokok?",
			Category:  "Perilaku Merokok",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			QuizID:    quizID,
			Text:      "Apa jenis rokok yang biasa dihisap?",
			Category:  "Perilaku Merokok",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	if err := db.Create(&smokingQuestions).Error; err != nil {
		return err
	}

	// Create answers for smoking behavior questions
	smokingStatusAnswers := []entity.Answer{
		{
			QuestionID: smokingQuestions[0].ID,
			Text:       "Iya, perokok aktif",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: smokingQuestions[0].ID,
			Text:       "Iya, perokok pasif",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: smokingQuestions[0].ID,
			Text:       "Tidak merokok",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&smokingStatusAnswers).Error; err != nil {
		return err
	}

	smokingStartAgeAnswers := []entity.Answer{
		{
			QuestionID: smokingQuestions[1].ID,
			Text:       "< 15 tahun",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: smokingQuestions[1].ID,
			Text:       "> 15 tahun",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: smokingQuestions[1].ID,
			Text:       "Tidak Merokok",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&smokingStartAgeAnswers).Error; err != nil {
		return err
	}

	smokingAmountAnswers := []entity.Answer{
		{
			QuestionID: smokingQuestions[2].ID,
			Text:       "Perokok ringan (< 10 batang perhari)",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: smokingQuestions[2].ID,
			Text:       "Perokok sedang (10-20 batang perhari)",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: smokingQuestions[2].ID,
			Text:       "Perokok berat (> 20 batang perhari)",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: smokingQuestions[2].ID,
			Text:       "Tidak Merokok",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&smokingAmountAnswers).Error; err != nil {
		return err
	}

	smokingDurationAnswers := []entity.Answer{
		{
			QuestionID: smokingQuestions[3].ID,
			Text:       "≤ 10 tahun",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: smokingQuestions[3].ID,
			Text:       "≥ 10 tahun",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: smokingQuestions[3].ID,
			Text:       "Tidak Merokok",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&smokingDurationAnswers).Error; err != nil {
		return err
	}

	smokingTypeAnswers := []entity.Answer{
		{
			QuestionID: smokingQuestions[4].ID,
			Text:       "Rokok Kretek",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: smokingQuestions[4].ID,
			Text:       "Rokok Putih",
			Value:      1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			QuestionID: smokingQuestions[4].ID,
			Text:       "Tidak Merokok",
			Value:      0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	if err := db.Create(&smokingTypeAnswers).Error; err != nil {
		return err
	}

	return nil
}
