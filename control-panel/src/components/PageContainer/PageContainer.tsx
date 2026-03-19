import type { ReactNode } from "react";
import Container from "react-bootstrap/Container";

type PageContainerProps = {
  title: string;
  children?: ReactNode;
};

function PageContainer({ title, children }: PageContainerProps) {
  return (
    <Container className="py-4">
      <h1 className="page-title">{title}</h1>
      {children}
    </Container>
  );
}

export default PageContainer;
