import Image from "next/image";

/** Show the complete original artwork; the viewport only trims blank canvas. */
export function BrandLogo({ priority = false }: { priority?: boolean }) {
  return <span className="wk-logo"><Image src="/brand/semiclaw-original.png" alt="SemiClaw" width={945} height={650} priority={priority} /></span>;
}
