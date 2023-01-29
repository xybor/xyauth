var userName = document.getElementById("user-name").value;

class headerComponent extends HTMLElement {
    connectedCallback() {
      this.innerHTML = `<header class="fixed-top">
      <nav class="navbar bg-light navbar-expand-sm navbar-toggleable-sm navbar-light">
          <div class="container">
              <a class="navbar-brand" href="/">
                  <img src="/static/img/logo_black.png" width="99" height="54" alt="logo"/>
              </a>
              <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#menu-bar, #sidebarMenu" aria-controls="navbarSupportedContent"
                      aria-expanded="false" aria-label="Toggle navigation">
                  <span class="navbar-toggler-icon"></span>
              </button>
              <div id="menu-bar" class="navbar-collapse collapse d-sm-inline-flex justify-content-end">
                  <ul class="navbar-nav">
                      <li class="nav-item mx-4">
                          <a class="nav-link text-dark" href="/" title="Home">Home</a>
                      </li>
                      <li class="nav-item mx-4">
                          <a class="nav-link text-dark" href="/register" title="Create Account">Create Account</a>
                      </li>
                      <li class="nav-item mx-4">
                          <div class="dropdown">
                              <a class="nav-link text-dark" data-toggle="dropdown" href="#" id="dropdownUser" data-bs-toggle="dropdown" aria-expanded="false" title="` + userName + `">` + userName +`</a>
                              <ul class="dropdown-menu" aria-labelledby="dropdownUser">
                                  <li>
                                      <a class="dropdown-item nav-link text-dark" href="/profile" title="Profile"><i class="fa-regular fa-user me-3"></i>Profile</a>
                                  </li>
                                  <li>
                                      <a class="dropdown-item nav-link text-dark" href="/logout?redirect_uri=/" title="Logout">Logout</a>
                                  </li>
                              </ul>
                          </div>


                      </li>
                      
                  </ul>
              </div>
          </div>
      </nav>
  </header>`;
    }
  }

customElements.define('header-component', headerComponent);

class footerComponent extends HTMLElement {
    connectedCallback() {
      this.innerHTML = `<footer class="d-flex align-items-center footer text-white bg-primary">
      <div class="container">
          <div class="row">
              <div class="col-md-6 col-sm-12 col-sm-12 text-center text-md-start">
                  <small>
                      <span>
                          &copy;
                          <script>
                              var CurrentYear = new Date().getFullYear();
                              document.write(CurrentYear);
                          </script>
                          - Xybor. All Rights Reserved.
                      </span>
                  </small>
                  <br />
                  <small>
                      <div>Follow us on:</div>
                      <a target="_blank" class="text-white" href="https://github.com/xybor" title="github" rel="noopener"><i class="fab fa-github fa-2x"></i></a>
                  </small>
              </div>
              <div class="col-md-6 col-sm-12 col-sm-12 text-center text-md-end">
                  <small>
                      <a target="_blank"
                         class="text-white text-decoration-none"
                         href="#">Privacy Policy</a>
                      |
                      <a target="_blank" class="text-white text-decoration-none" href="#">Terms & Conditions</a>
                  </small>
                  <br />
                  <img src="/static/img/logo_white.png" width="99" height="54" alt="logo"/>
              </div>
          </div>
      </div>
  </footer>`;
    }
  }
      
customElements.define('footer-component', footerComponent);

class sidebarComponent extends HTMLElement {
    connectedCallback() {
      this.innerHTML = `<div class="position-sticky pt-3">
        <ul class="nav flex-column">
          <li class="nav-item">
            <a class="nav-link text-dark" aria-current="page" href="#">
              <i class="fa-regular fa-user me-3"></i>
              User
            </a>
          </li>
          <li class="nav-item">
            <a class="nav-link text-dark" href="#">
              <i class="fa-regular fa-clock me-3"></i>
              Sessions
            </a>
          </li>
        </ul>
      </div>`;
    }
  }
      
customElements.define('sidebar-component', sidebarComponent);
