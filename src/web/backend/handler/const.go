package handler

const ( // Paths to templates
	homePageTemplatePath             = "./src/web/frontend/templates/home.page.tmpl"
	navbarTemplatePath               = "./src/web/frontend/templates/navbar.tmpl"
	basePageTemplatePath             = "./src/web/frontend/templates/base.page.tmpl"
	portfolioCreatorPageTemplatePath = "./src/web/frontend/templates/portfolio_creator.page.tmpl"
	portfolioChooserPageTemplatePath = "./src/web/frontend/templates/portfolio_chooser.page.tmpl"
	managerPageTemplatePath          = "./src/web/frontend/templates/manager.page.tmpl"
	followingIndexPageTemplatePath   = "./src/web/frontend/templates/following_index.page.tmpl"
	followingIndexTableTemplatePath  = "./src/web/frontend/templates/following_index.table.tmpl"
)

const ( // Cookie names
	portfolioCookieName = "current_portfolio"
)

const ( // Data_to_render keys
	allPortfoliosKey = "portfolios"
	allAssetsKey     = "assets"
	allIndexesKey    = "indexes"
	pieChartKey      = "PieChart"
)

const ( // Keys for forms from frontend
	indexFormKey     = "index"
	portfolioFormKey = "portfolio-name"
	allAssetsFormKey = "assets[]"
)

const ( // Others
	frontendFolderPath = "src/web/frontend"
)
