package server

import (
	"net/http"

	"github.com/err0r500/go-realworld-clean/domain"
	"github.com/err0r500/go-realworld-clean/implem/json.formatter"
	"github.com/gin-gonic/gin"
)

type ArticleReq struct {
	Article struct {
		Title       string   `json:"title,required"`
		Description string   `json:"description,required"`
		Body        string   `json:"body,required"`
		TagList     []string `json:"tagList,required"`
	} `json:"article,required"`
}

func articleFromReq(req *ArticleReq) domain.Article {
	return domain.Article{
		Title:       req.Article.Title,
		Description: req.Article.Description,
		Body:        req.Article.Body,
		TagList:     req.Article.TagList,
	}
}

func (rH RouterHandler) articlePost(c *gin.Context) {
	log := rH.log(c.Request.URL.Path)

	userName, err := rH.getUserName(c)
	if err != nil {
		log(err)
		c.Status(http.StatusUnauthorized)
		return
	}

	req := &ArticleReq{}
	if err := c.BindJSON(req); err != nil {
		log(err)
		c.Status(http.StatusBadRequest)
		return
	}

	article, err := rH.ucHandler.ArticlePost(userName, articleFromReq(req))
	if err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"article": formatter.NewArticleFromDomain(*article, true)})
}

func (rH RouterHandler) articlePut(c *gin.Context) {
	log := rH.log(c.Request.URL.Path)
	userName, err := rH.getUserName(c)
	if err != nil {
		log(err)
		c.Status(http.StatusUnauthorized)
		return
	}

	req := &ArticleReq{}
	if err := c.BindJSON(req); err != nil {
		log(err)
		c.Status(http.StatusBadRequest)
		return
	}
	article, err := rH.ucHandler.ArticlePut(userName, c.Param("slug"), articleFromReq(req))
	if err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}

	c.JSON(http.StatusOK, gin.H{"article": formatter.NewArticleFromDomain(*article, true)})
}

func (rH RouterHandler) articleGet(c *gin.Context) {
	log := rH.log(c.Request.URL.Path)

	article, err := rH.ucHandler.ArticleGet(c.Param("slug"))
	if err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}
	c.JSON(http.StatusOK, gin.H{"article": formatter.NewArticleFromDomain(*article, true)})
}

func (rH RouterHandler) articleDelete(c *gin.Context) {
	log := rH.log(c.Request.URL.Path)

	userName, err := rH.getUserName(c)
	if err != nil {
		log(err)
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := rH.ucHandler.ArticleDelete(userName, c.Param("slug")); err != nil {
		log(err)
		c.Status(http.StatusUnprocessableEntity)
		return
	}
	c.Status(http.StatusOK)
}
