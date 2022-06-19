class PagesController < ApplicationController
  before_action :set_page, only: %i[ show edit update destroy ]

  def page_by_slug_and_domain
    @page = Page.joins({:website => :domains}).where(website: {domains: {name: request.host}}, pages: {slug: params[:path]}).first!
    # https://stackoverflow.com/questions/9560093/bypass-application-html-erb-in-rails-app
    render "page_from_db", :layout => false
  end

  def page_by_slug
    @page = Page.joins(:website).where(website: {slug: params[:website]}, pages: {slug: params[:path]}).first!
    # render :layout => false
    render "page_from_db", :layout => false
  end

  # GET /pages or /pages.json
  def index
    @pages = Page.all
  end

  # GET /pages/1 or /pages/1.json
  def show
  end

  # GET /pages/new
  def new
    @page = Page.new
  end

  # GET /pages/1/edit
  def edit
  end

  # POST /pages or /pages.json
  def create
    @page = Page.new(page_params)

    respond_to do |format|
      if @page.save
        format.html { redirect_to page_url(@page), notice: "Page was successfully created." }
        format.json { render :show, status: :created, location: @page }
      else
        format.html { render :new, status: :unprocessable_entity }
        format.json { render json: @page.errors, status: :unprocessable_entity }
      end
    end
  end

  # PATCH/PUT /pages/1 or /pages/1.json
  def update
    respond_to do |format|
      if @page.update(page_params)
        format.html { redirect_to page_url(@page), notice: "Page was successfully updated." }
        format.json { render :show, status: :ok, location: @page }
      else
        format.html { render :edit, status: :unprocessable_entity }
        format.json { render json: @page.errors, status: :unprocessable_entity }
      end
    end
  end

  # DELETE /pages/1 or /pages/1.json
  def destroy
    @page.destroy

    respond_to do |format|
      format.html { redirect_to pages_url, notice: "Page was successfully destroyed." }
      format.json { head :no_content }
    end
  end

  private
    # Use callbacks to share common setup or constraints between actions.
    def set_page
      @page = Page.find(params[:id])
    end

    # Only allow a list of trusted parameters through.
    def page_params
      params.require(:page).permit(:website_id, :user_id, :breed, :language, :active_from, :slug, :title, :body)
    end
end
