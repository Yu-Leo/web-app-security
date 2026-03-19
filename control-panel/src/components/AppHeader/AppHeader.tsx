import Container from "react-bootstrap/Container";
import Nav from "react-bootstrap/Nav";
import Navbar from "react-bootstrap/Navbar";
import { NavLink } from "react-router-dom";
import logoMark from "../../assets/brand/was-logo-mark.png";
import Icon from "../Icon";

function AppHeader() {
  return (
    <Navbar bg="light" expand="lg" className="border-bottom app-navbar">
      <Container>
        <Navbar.Brand
          as={NavLink}
          to="/resources"
          className="d-flex align-items-center gap-2"
        >
          <img src={logoMark} alt="WAS" className="brand-logo" />
          <span>Web App Security</span>
        </Navbar.Brand>
        <Navbar.Toggle aria-controls="main-nav" />
        <Navbar.Collapse id="main-nav">
          <Nav className="me-auto">
            <Nav.Link as={NavLink} to="/resources">
              <Icon name="diagram-3" className="me-1" />
              Ресурсы
            </Nav.Link>
            <Nav.Link as={NavLink} to="/security-profiles">
              <Icon name="shield-check" className="me-1" />
              Профили безопасности
            </Nav.Link>
            <Nav.Link as={NavLink} to="/traffic-profiles">
              <Icon name="speedometer2" className="me-1" />
              Профили ограничителя трафика
            </Nav.Link>
            <Nav.Link as={NavLink} to="/event-logs">
              <Icon name="bell" className="me-1" />
              Логи событий
            </Nav.Link>
            <Nav.Link as={NavLink} to="/request-logs">
              <Icon name="list-ul" className="me-1" />
              Логи запросов
            </Nav.Link>
            <Nav.Link as={NavLink} to="/ml-models">
              <Icon name="diagram-3" className="me-1" />
              ML-модели
            </Nav.Link>
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
}

export default AppHeader;
