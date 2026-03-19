import Pagination from "react-bootstrap/Pagination";
import Form from "react-bootstrap/Form";

type TablePaginationProps = {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  pageSize?: number;
  pageSizeOptions?: number[];
  onPageSizeChange?: (size: number) => void;
};

function TablePagination({
  currentPage,
  totalPages,
  onPageChange,
  pageSize,
  pageSizeOptions = [10, 20, 50],
  onPageSizeChange,
}: TablePaginationProps) {
  const pages: number[] = [];
  const safeTotalPages = Math.max(1, totalPages);
  const safeCurrentPage = Math.min(Math.max(1, currentPage), safeTotalPages);
  const start = Math.max(1, safeCurrentPage - 2);
  const end = Math.min(safeTotalPages, safeCurrentPage + 2);

  for (let page = start; page <= end; page += 1) {
    pages.push(page);
  }

  return (
    <div className="mt-3 border rounded bg-body-tertiary px-4 py-3">
      <div className="d-flex flex-wrap justify-content-center align-items-center gap-3">
        {pageSize && onPageSizeChange ? (
          <div className="d-flex align-items-center gap-2">
            <span className="text-muted small">Элементов на странице</span>
            <Form.Select
              className="w-auto"
              value={pageSize}
              onChange={(event) => onPageSizeChange(Number(event.target.value))}
            >
              {pageSizeOptions.map((option) => (
                <option key={option} value={option}>
                  {option}
                </option>
              ))}
            </Form.Select>
          </div>
        ) : null}
        <Pagination className="mb-0">
          <Pagination.First
            onClick={() => onPageChange(1)}
            disabled={safeCurrentPage === 1}
          />
          <Pagination.Prev
            onClick={() => onPageChange(safeCurrentPage - 1)}
            disabled={safeCurrentPage === 1}
          />
          {start > 1 ? <Pagination.Ellipsis disabled /> : null}
          {pages.map((page) => (
            <Pagination.Item
              key={page}
              active={page === safeCurrentPage}
              onClick={() => onPageChange(page)}
            >
              {page}
            </Pagination.Item>
          ))}
          {end < safeTotalPages ? <Pagination.Ellipsis disabled /> : null}
          <Pagination.Next
            onClick={() => onPageChange(safeCurrentPage + 1)}
            disabled={safeCurrentPage === safeTotalPages}
          />
          <Pagination.Last
            onClick={() => onPageChange(safeTotalPages)}
            disabled={safeCurrentPage === safeTotalPages}
          />
        </Pagination>
      </div>
    </div>
  );
}

export default TablePagination;
