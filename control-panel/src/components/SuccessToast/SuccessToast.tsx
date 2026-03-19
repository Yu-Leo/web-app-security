import Toast from "react-bootstrap/Toast";
import ToastContainer from "react-bootstrap/ToastContainer";

type SuccessToastProps = {
  show: boolean;
  onClose: () => void;
};

function SuccessToast({ show, onClose }: SuccessToastProps) {
  return (
    <ToastContainer position="top-end" className="p-3">
      <Toast show={show} onClose={onClose} delay={2500} autohide bg="secondary">
        <Toast.Body className="text-white">Изменения сохранены</Toast.Body>
      </Toast>
    </ToastContainer>
  );
}

export default SuccessToast;
