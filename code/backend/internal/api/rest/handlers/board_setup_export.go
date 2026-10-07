package handlers

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"digital-innovation/gostrategy/internal/api/core"
	"digital-innovation/gostrategy/internal/db"
	"digital-innovation/gostrategy/internal/game"

	"github.com/gin-gonic/gin"
)

// ExportBoardSetupHandler exports a single board setup as JSON for external AI training/testing
// @Summary Export board setup
// @Description Export a board setup configuration in JSON format for external AI training/testing
// @Tags board-setups
// @Produce json
// @Param id path int true "Board Setup ID"
// @Success 200 {object} game.ExportedBoardSetup
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Board setup not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /board-setups/{id}/export [get]
func (h *Handler) ExportBoardSetupHandler(c *gin.Context) {
	user := core.EnsureAuthenticated(c)
	if user == nil {
		return
	}

	id, err := core.ParseID(c, "id")
	if err != nil || id <= 0 {
		core.SendError(c, "Invalid board setup ID", http.StatusBadRequest)
		return
	}
	setup, err := db.GetBoardSetup(c.Request.Context(), int(id), user.ID)
	if err != nil {
		core.SendError(c, "Board setup not found", http.StatusNotFound)
		return
	}

	rows, err := game.SetupDataToRows(setup.SetupData)
	if err != nil {
		rows = []string{}
	}

	exported := game.ExportedBoardSetup{
		Name:        setup.Name,
		Description: setup.Description,
		SetupData:   setup.SetupData,
		Rows:        rows,
	}

	filename := sanitizeExportFilename(setup.Name)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.json\"", filename))
	core.SendJSON(c, exported, http.StatusOK)
}

// ExportAllBoardSetupsHandler exports all board setups as a ZIP archive of individual JSON files
// @Summary Export all board setups
// @Description Export all saved board setups as a ZIP archive containing individual JSON files
// @Tags board-setups
// @Produce application/zip
// @Success 200 {file} binary
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /board-setups/export [get]
func (h *Handler) ExportAllBoardSetupsHandler(c *gin.Context) {
	user := core.EnsureAuthenticated(c)
	if user == nil {
		return
	}

	setups, err := db.GetUserBoardSetups(c.Request.Context(), user.ID)
	if err != nil {
		core.SendError(c, "Failed to retrieve board setups", http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	usedNames := make(map[string]bool)
	for _, s := range setups {
		rows, err := game.SetupDataToRows(s.SetupData)
		if err != nil {
			rows = []string{}
		}

		item := game.ExportedBoardSetup{
			Name:        s.Name,
			Description: s.Description,
			SetupData:   s.SetupData,
			Rows:        rows,
		}

		jsonData, err := json.MarshalIndent(item, "", "  ")
		if err != nil {
			core.SendError(c, "Failed to serialize board setup", http.StatusInternalServerError)
			return
		}

		baseName := sanitizeExportFilename(s.Name)
		fileName := baseName + ".json"
		counter := 2
		for usedNames[fileName] {
			fileName = fmt.Sprintf("%s_%d.json", baseName, counter)
			counter++
		}
		usedNames[fileName] = true

		w, err := zipWriter.Create(fileName)
		if err != nil {
			core.SendError(c, "Failed to create zip entry", http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(jsonData); err != nil {
			core.SendError(c, "Failed to write zip content", http.StatusInternalServerError)
			return
		}
	}

	if err := zipWriter.Close(); err != nil {
		core.SendError(c, "Failed to generate zip file", http.StatusInternalServerError)
		return
	}

	c.Header("Content-Disposition", "attachment; filename=\"board_setups.zip\"")
	c.Data(http.StatusOK, "application/zip", buf.Bytes())
}

func sanitizeExportFilename(name string) string {
	cleaned := strings.ToLower(strings.TrimSpace(name))
	var sb strings.Builder
	for _, r := range cleaned {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			sb.WriteRune(r)
		} else if r == ' ' {
			sb.WriteRune('_')
		}
	}
	res := sb.String()
	if res == "" {
		return "board_setup"
	}
	return res
}
