import {Tag} from "@arco-design/web-react";

export function makeNewSearch(keyword: string, router) {
  // change search input in navbar
  const searchInput = document.getElementById('search-input');
  if (searchInput) {
    searchInput.setAttribute('value', keyword);
  }
  router.push({
    pathname: '/search',
    query: {
      q: keyword,
    },
  });
}

export const parseKeyword = (keyword: string, router) => {
  // if keyword does not have space, return it directly
  if (keyword === undefined) {
    return null;
  }
  if (keyword.indexOf(' ') === -1) {
    return null;
  }
  const keywords = keyword.split(' ');
  // do not split them to multiple div element
  return keywords.map((keyword, index) => (
    <Tag
      key={index.toString()}
      onClick={(event) => {
        makeNewSearch(keyword, router);
        event.stopPropagation();
      }}
      style={{
        cursor: 'pointer',
        marginRight: '4px',
        marginBottom: '4px',
        backgroundColor: 'rgba(var(--gray-6), 0.4)',
      }}
    >
      {keyword}
    </Tag>
  ));
};

export const parseKeywordOnVideo = (keyword: string, router) => {
  // if keyword does not have space, return it directly
  if (keyword === undefined) {
    return null;
  }
  if (keyword.indexOf(' ') === -1) {
    return null;
  }
  const keywords = keyword.split(' ');
  // do not split them to multiple div element
  return keywords.map((keyword, index) => (
    <Tag
      key={index.toString()}
      onClick={(event) => {
        makeNewSearch(keyword, router);
        event.stopPropagation();
      }}
      style={{
        cursor: 'pointer',
        marginRight: '4px',
        marginBottom: '4px',
      }}
    >
      {keyword}
    </Tag>
  ));
};