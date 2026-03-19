import Alert from "react-bootstrap/Alert";
import Icon from "../Icon";

type ErrorAlertProps = {
  error: string | null;
};

function ErrorAlert({ error }: ErrorAlertProps) {
  if (!error) {
    return null;
  }
  return (
    <Alert variant="danger" className="mb-3">
      <Icon name="exclamation-triangle" className="me-2" />
      {error}
    </Alert>
  );
}

export default ErrorAlert;
